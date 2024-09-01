package internal

import (
    "log"
    "os"

    "github.com/redis/go-redis/v9"
)

var (
    rdb           *redis.Client
    redisActive   bool
)

func InitRedis() {
    if redisActive {
        redisAddr := os.Getenv("REDIS_ADDR")
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
