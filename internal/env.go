package internal

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
    if err := godotenv.Load(); err != nil {
        log.Fatalf("Error loading .env file")
    }
    redisActive = os.Getenv("REDIS_ACTIVE") == "true"
}
