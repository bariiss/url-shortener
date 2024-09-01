package main

import (
	i "github.com/bariiss/url-shortener/internal"
)

func main() {
    i.LoadEnv()
    i.InitRedis()

    app := i.InitFiberApp()

    i.StartServer(app)
}
