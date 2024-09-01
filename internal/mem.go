package internal

import (
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

// InitMem initializes the memory storage for the URL shortener service.
//
// If the Redis active flag is set, it creates a new Redis client and pings the
// Redis server to ensure a connection can be established. If the connection
// fails, the function logs an error message and exits the program.
//
// If the Redis active flag is not set, it creates an empty map to store the
// short URLs and logs a message indicating that in-memory storage is being
// used.
func InitMem() {
	if redisActive {
		rdb = redis.NewClient(&redis.Options{
			Addr: redisAddr,
		})

		if _, err := rdb.Ping(ctx).Result(); err != nil {
			log.Fatalf("Error connecting to Redis: %v", err)
		}
		log.Println("Connected to Redis")
		return
	}

	memoryStorage = make(map[string]string)
	log.Println("Redis is not active. Using in-memory storage.")
}

// setURL sets a short URL to the given original URL in the underlying store.
//
// If Redis is active, it will check if the short URL already exists in Redis.
// If it does, it will check if the existing URL matches the given original URL.
// If it does, it will return nil. If it does not, it will generate a new short
// URL and repeat the process until a unique short URL is found.
//
// If Redis is not active, it will use an in-memory map to store the URLs. It
// will use a lock to ensure thread safety.
//
// The function will return an error if there is an issue with the Redis
// connection or if there is an error generating a new short URL.
func setURL(shortURL, originalURL string) error {
	if redisActive {
		for {
			existingURL, err := rdb.Get(ctx, shortURL).Result()
			if err == redis.Nil {
				return rdb.Set(ctx, shortURL, originalURL, 0).Err()
			}

			if err != nil {
				return fmt.Errorf("error checking Redis: %v", err)
			}

			if existingURL == originalURL {
				return nil
			}

			shortURL = generateShortURL()
		}
	}

	for {
		mu.RLock()
		existingURL, exists := memoryStorage[shortURL]
		mu.RUnlock()

		if !exists {
			mu.Lock()
			memoryStorage[shortURL] = originalURL
			mu.Unlock()
			return nil
		}

		if existingURL == originalURL {
			return nil
		}

		shortURL = generateShortURL()
	}
}

// getURL retrieves the original URL associated with the given short URL.
//
// If Redis is active, it will retrieve the URL from Redis. Otherwise, it will
// retrieve the URL from an in-memory map. If the URL is not found, it will
// return an error.
//
// The function is thread-safe when Redis is inactive.
func getURL(shortURL string) (string, error) {
	if redisActive {
		return rdb.Get(ctx, shortURL).Result()
	}
	mu.RLock()
	defer mu.RUnlock()
	originalURL, exists := memoryStorage[shortURL]
	if !exists {
		return "", fmt.Errorf("URL not found")
	}
	return originalURL, nil
}
