package internal

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
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

func indexHandler(c *fiber.Ctx) error {
    return c.Render("index", nil)
}
