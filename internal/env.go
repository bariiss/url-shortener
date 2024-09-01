package internal

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	redisActive bool
	letterBytes string
	appPort     string
	redisAddr   string
	maxRequests int
	expiration  int
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}

	redisActive = os.Getenv("REDIS_ACTIVE") == "true"
	letterBytes = os.Getenv("LETTER_BYTES")
	appPort = os.Getenv("APP_PORT")
	redisAddr = os.Getenv("REDIS_ADDR")

	maxRequests, err := strconv.Atoi(os.Getenv("MAX_REQUESTS"))
	if err != nil {
		log.Fatalf("Invalid MAX_REQUESTS value: %v", err)
	}

	expiration, err = strconv.Atoi(os.Getenv("EXPIRATION"))
	if err != nil {
		log.Fatalf("Invalid EXPIRATION value: %v", err)
	}

	log.Println("----------------------")
	log.Println("Environment Variables")
	log.Println("----------------------")
	log.Println("REDIS_ACTIVE:", redisActive)
	log.Println("APP_PORT:", appPort)
	log.Println("REDIS_ADDR:", redisAddr)
	log.Println("MAX_REQUESTS:", maxRequests)
	log.Println("EXPIRATION:", expiration)
	log.Println("----------------------")
}
