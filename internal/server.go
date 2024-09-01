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

type AppConfig struct {
	Port   string
	Engine *html.Engine
}

func SetAppConfig() *AppConfig {
	return &AppConfig{
		Port:   appPort,
		Engine: initTemplateEngine(),
	}
}

func StartServer(config *AppConfig) {
	showWelcomeMessage()
	app := initFiberApp(config)
	log.Printf("Starting server on :%s", config.Port)
	if err := app.Listen(fmt.Sprintf(":%s", config.Port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func showWelcomeMessage() {
    fmt.Println(`
    ====================================
    Welcome to the URL Shortener Service
    ====================================
    `)
}

func initTemplateEngine() *html.Engine {
    return html.New("./templates", ".html")
}

func initFiberApp(config *AppConfig) *fiber.App {
	app := fiber.New(fiber.Config{
		Views: config.Engine,
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

func shortenHandler(c *fiber.Ctx) error {
    originalURL := c.FormValue("url")
    shortURL := generateShortURL()

    err := setURL(shortURL, originalURL)
    if err != nil {
        log.Printf("Error storing URL: %v", err)
        return c.Status(fiber.StatusInternalServerError).SendString("Error storing URL")
    }

    fullShortURL := fmt.Sprintf("%s/r/%s", c.BaseURL(), shortURL)
    log.Printf("Stored URL: %s -> %s", shortURL, originalURL)
    response := fmt.Sprintf(`
        <p>Shortened URL: <a href="%s" target="_blank">%s</a>
        <span class="copy-icon" onclick="copyToClipboard('%s')">📋</span></p>
    `, fullShortURL, fullShortURL, fullShortURL)
    return c.SendString(response)
}

func redirectHandler(c *fiber.Ctx) error {
    shortURL := c.Params("shortURL")
    originalURL, err := getURL(shortURL)
    if err != nil {
        log.Printf("Error retrieving URL for %s: %v", shortURL, err)
        return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving URL")
    }
    if originalURL == "" {
        log.Printf("URL not found for %s", shortURL)
        return c.Status(fiber.StatusNotFound).SendString("URL not found")
    }
    log.Printf("Redirecting %s to %s", shortURL, originalURL)
    return c.Redirect(originalURL, fiber.StatusSeeOther)
}
