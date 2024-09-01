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
	Expiration  time.Duration
}

// SetAppConfig returns a new AppConfig instance with the given configuration.
//
// It takes no arguments and returns a pointer to an AppConfig instance.
//
// The returned AppConfig instance is populated with the following values:
//
// - Port: The value of appPort.
// - Engine: The result of calling initTemplateEngine().
// - MaxRequests: The value of maxRequests.
// - Expiration: The value of expiration converted to a time.Duration.
func SetAppConfig() *AppConfig {
	return &AppConfig{
		Port:        appPort,
		Engine:      initTemplateEngine(),
		MaxRequests: maxRequests,
		Expiration:  time.Duration(expiration),
	}
}

// StartServer starts the Fiber server with the given configuration.
//
// It prints a welcome message to the console, initializes the Fiber app,
// and starts the server on the configured port. If the server fails to start,
// it logs an error message and exits the program.
func StartServer(config *AppConfig) {
	showWelcomeMessage()
	app := initFiberApp(config)
	log.Printf("Starting server on :%s", config.Port)
	if err := app.Listen(fmt.Sprintf(":%s", config.Port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// showWelcomeMessage prints a welcome message to the console when the server is started.
func showWelcomeMessage() {
	fmt.Println(`
    ====================================
    Welcome to the URL Shortener Service
    ====================================
    `)
}

// initTemplateEngine initializes a new html template engine using the ./templates directory and .html as the file extension.
func initTemplateEngine() *html.Engine {
	return html.New("./templates", ".html")
}

// initFiberApp initializes a new Fiber app with the given configuration.
//
// It sets up the routing for the URL shortener service:
// - POST /shorten: handles URL shortening requests with rate limiting
// - GET /: renders the index.html template
// - GET /r/:shortURL: redirects to the stored URL
// - /static: serves the static files from the ./static directory
//
// The rate limiting is based on the client's IP address, with a maximum number
// of requests allowed within a given expiration time. If the limit is reached,
// a 429 status code is returned with an HTML message indicating the error.
func initFiberApp(config *AppConfig) *fiber.App {
	app := fiber.New(fiber.Config{
		Views:       config.Engine,
		ProxyHeader: fiber.HeaderXForwardedFor,
	})

	app.Post("/shorten", limiter.New(limiter.Config{
		Max:        config.MaxRequests,
		Expiration: config.Expiration * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return getClientIP(c)
		},
		LimitReached: func(c *fiber.Ctx) error {
			clientIP := getClientIP(c)
			message := fmt.Sprintf(`Too many requests from IP: <strong>%s</strong>.`, clientIP)
			return c.Status(fiber.StatusTooManyRequests).SendString(message)
		},
	}), shortenHandler)

	app.Get("/", indexHandler)
	app.Get("/r/:shortURL", redirectHandler)
	app.Static("/static", "./static")

	return app
}

// indexHandler handles GET requests to / and renders the index.html template.
func indexHandler(c *fiber.Ctx) error {
	return c.Render("index", nil)
}

// shortenHandler handles POST requests to /shorten and stores the original URL
// in the Redis store. If the URL is not found, it returns a 404 status code.
// If there is an error retrieving the URL, it returns a 500 status code.
//
// The body of the request should contain the URL to shorten as a form field
// named "url".
//
// The response is a string containing the shortened URL with a copy icon.
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

// redirectHandler handles GET requests to /r/:shortURL and redirects to the stored URL.
// If the URL is not found, it returns a 404 status code.
// If there is an error retrieving the URL, it returns a 500 status code.
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

// getClientIP returns the client's IP address.
//
// If the request contains the X-Forwarded-For header, it splits the value by
// comma and returns the first element. Otherwise, it returns the IP address
// from the request context.
func getClientIP(c *fiber.Ctx) string {
	clientIP := c.Get("X-Forwarded-For")
	if clientIP != "" {
		return strings.Split(clientIP, ",")[0]
	}

	return c.IP()
}
