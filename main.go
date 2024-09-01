package main

import (
	i "github.com/bariiss/url-shortener/internal"
)

func init() {
	i.LoadEnv()
}

func main() {
    i.InitMem()

    config := i.SetAppConfig()
    i.StartServer(config)
}
