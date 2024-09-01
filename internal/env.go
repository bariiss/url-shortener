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

// LoadEnv loads environment variables from the .env file.
//
// The following environment variables are supported:
//
// - REDIS_ACTIVE: Set to "true" to enable Redis as the store.
// - LETTER_BYTES: The characters to use when generating short URLs.
// - APP_PORT: The port number to listen on.
// - REDIS_ADDR: The address of the Redis server.
// - MAX_REQUESTS: The maximum number of requests allowed within a given expiration time.
// - EXPIRATION: The expiration time in seconds for the request limit.
//
// If any of the variables are invalid, the function will log an error message and exit the program.
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
