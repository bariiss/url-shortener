package internal

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/template/html/v2"
)

var (
	engine *html.Engine
)

func init() {
	engine = initTemplateEngine()
}

func showWelcomeMessage() {
    fmt.Println(`
    ====================================
    Welcome to the URL Shortener Service
    ====================================
    `)
}

func StartServer(app *fiber.App) {
	showWelcomeMessage()
    log.Printf("Starting server on :%s", appPort)
    log.Fatal(app.Listen(fmt.Sprintf(":%s", appPort)))
}

func initTemplateEngine() *html.Engine {
    return html.New("./templates", ".html")
}

func InitFiberApp() *fiber.App {
    app := fiber.New(fiber.Config{
        Views: engine,
    })

    app.Post("/shorten", limiter.New(limiter.Config{
        Max:        5,
        Expiration: 60 * time.Second,
        KeyGenerator: func(c *fiber.Ctx) string {
            return c.IP()
        },
        LimitReached: func(c *fiber.Ctx) error {
            return c.Status(fiber.StatusTooManyRequests).SendString("<p>Too many requests. Please try again later.</p>")
        },
    }), shortenHandler)

    app.Get("/", indexHandler)
    app.Get("/r/:shortURL", redirectHandler)
    app.Static("/static", "./static")

    return app
}
