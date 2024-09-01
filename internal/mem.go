package internal

import (
	"fmt"
    "log"

    "github.com/redis/go-redis/v9"
)

var (
    rdb           *redis.Client
)

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
