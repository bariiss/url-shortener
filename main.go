package main

func main() {
	loadEnv()
	initRedis()

	engine := initTemplateEngine()
	app := initFiberApp(engine)

	startServer(app)
}
