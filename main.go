package main

func main() {
    showWelcomeMessage()
    loadEnv()
    initRedis()

    engine := initTemplateEngine()
    app := initFiberApp(engine)

    startServer(app)
}
