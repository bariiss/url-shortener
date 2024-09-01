package main

import (
	i "github.com/bariiss/url-shortener/internal"
)

func main() {
    i.ShowWelcomeMessage()
    i.LoadEnv()
    i.InitRedis()

    engine := i.InitTemplateEngine()
    app := i.InitFiberApp(engine)

    i.StartServer(app)
}
