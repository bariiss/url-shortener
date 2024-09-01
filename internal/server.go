package internal

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func ShowWelcomeMessage() {
    fmt.Println(`
    ====================================
    Welcome to the URL Shortener Service
    ====================================
    `)
}

func StartServer(app *fiber.App) {
    appPort := os.Getenv("APP_PORT")
    log.Printf("Starting server on :%s", appPort)
    log.Fatal(app.Listen(fmt.Sprintf(":%s", appPort)))
}

func InitTemplateEngine() *html.Engine {
    return html.New("./templates", ".html")
}

func InitFiberApp(engine *html.Engine) *fiber.App {
    app := fiber.New(fiber.Config{
        Views: engine,
    })

    app.Get("/", indexHandler)
    app.Post("/shorten", shortenHandler)
    app.Get("/r/:shortURL", redirectHandler)
    app.Static("/static", "./static")

    return app
}

func indexHandler(c *fiber.Ctx) error {
    return c.Render("index", nil)
}
