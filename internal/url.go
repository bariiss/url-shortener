package internal

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/exp/rand"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var (
	mu            sync.RWMutex
	memoryStorage map[string]string
	ctx           = context.Background()
)

func generateShortURL() string {
    rand.Seed(uint64(time.Now().UnixNano()))
    b := make([]byte, 8)
    for i := range b {
        b[i] = letterBytes[rand.Intn(len(letterBytes))]
    }
    return string(b)
}

func setURL(shortURL, originalURL string) error {
    if redisActive {
        return rdb.Set(ctx, shortURL, originalURL, 0).Err()
    } else {
        mu.Lock()
        defer mu.Unlock()
        memoryStorage[shortURL] = originalURL
        return nil
    }
}

func getURL(shortURL string) (string, error) {
    if redisActive {
        return rdb.Get(ctx, shortURL).Result()
    } else {
        mu.RLock()
        defer mu.RUnlock()
        originalURL, exists := memoryStorage[shortURL]
        if !exists {
            return "", fmt.Errorf("URL not found")
        }
        return originalURL, nil
    }
}
