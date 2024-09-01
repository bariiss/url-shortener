package main

import (
	i "github.com/bariiss/url-shortener/internal"
)

// main is the entry point of the URL shortener service.
//
// It loads the environment variables from the .env file, initializes the
// memory storage, and starts the Fiber server with the given configuration.
func main() {
	i.LoadEnv()
	i.InitMem()

	config := i.SetAppConfig()
	i.StartServer(config)
}
