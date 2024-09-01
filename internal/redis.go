package internal

import (
    "log"

    "github.com/redis/go-redis/v9"
)

var (
    rdb           *redis.Client
)

func InitRedis() {
    if redisActive {
        rdb = redis.NewClient(&redis.Options{
            Addr: redisAddr,
        })

        if _, err := rdb.Ping(ctx).Result(); err != nil {
            log.Fatalf("Error connecting to Redis: %v", err)
        }
        log.Println("Connected to Redis")
    } else {
        memoryStorage = make(map[string]string)
        log.Println("Redis is not active. Using in-memory storage.")
    }
}
