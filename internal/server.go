package internal

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/template/html/v2"
)

type AppConfig struct {
	Port        string
	Engine      *html.Engine
	MaxRequests int
	Expiration  int
}

var (
	engine *html.Engine
)

func SetAppConfig() *AppConfig {
	return &AppConfig{
		Port:        appPort,
		Engine:      initTemplateEngine(),
		MaxRequests: maxRequests,
		Expiration: expiration,
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
		ProxyHeader: fiber.HeaderXForwardedFor,
	})

	app.Post("/shorten", limiter.New(limiter.Config{
		Max:        config.MaxRequests,
		Expiration: 60 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return getClientIP(c)
		},
		LimitReached: func(c *fiber.Ctx) error {
			clientIP := getClientIP(c)
			message := fmt.Sprintf(`
				<div>Too many requests from IP: <strong>%s</strong>. Please try again later.</div>
			`, clientIP)
			return c.Status(fiber.StatusTooManyRequests).SendString(message)
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
	if originalURL == "" {
		return c.Status(fiber.StatusBadRequest).SendString("URL cannot be empty")
	}

	if !strings.HasPrefix(originalURL, "http://") && !strings.HasPrefix(originalURL, "https://") {
		originalURL = "https://" + originalURL
	}

	shortURL := generateShortURL()

	err := setURL(shortURL, originalURL)
	if err != nil {
		log.Printf("Error storing URL: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error storing URL")
	}

	fullShortURL := fmt.Sprintf("%s/r/%s", c.BaseURL(), shortURL)
	log.Printf("Stored URL: %s -> %s", shortURL, originalURL)
	response := fmt.Sprintf(`
        Shortened URL: <a href="%s" target="_blank">%s</a>
        <span class="copy-icon" onclick="copyToClipboard('%s')">📋</span>
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

func getClientIP(c *fiber.Ctx) string {
	clientIP := c.Get("X-Forwarded-For")
	if clientIP != "" {
		return strings.Split(clientIP, ",")[0]
	}

	return c.IP()
}
