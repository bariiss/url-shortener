package internal

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
    redisActive   bool
    letterBytes string
    appPort 	string
    redisAddr 	string
)

func LoadEnv() {
    if err := godotenv.Load(); err != nil {
        log.Fatalf("Error loading .env file")
    }
    redisActive = os.Getenv("REDIS_ACTIVE") == "true"
    letterBytes = os.Getenv("LETTER_BYTES")
    appPort = os.Getenv("APP_PORT")
    redisAddr = os.Getenv("REDIS_ADDR")
}
