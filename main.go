package main

import (
	i "github.com/bariiss/url-shortener/internal"
)

func main() {
	i.LoadEnv()
    i.InitMem()

    config := i.SetAppConfig()
    i.StartServer(config)
}
