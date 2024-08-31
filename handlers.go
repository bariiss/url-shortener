package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/rand"
)

var (
	rdb *redis.Client
	ctx = context.Background()
	mu  sync.RWMutex
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func initRedis() {
	redisAddr := os.Getenv("REDIS_ADDR")
	rdb = redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}
}

func initTemplateEngine() *html.Engine {
	return html.New("./templates", ".html")
}

func initFiberApp(engine *html.Engine) *fiber.App {
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", indexHandler)
	app.Post("/shorten", shortenHandler)
	app.Get("/r/:shortURL", redirectHandler)
	app.Static("/static", "./static")

	return app
}

func startServer(app *fiber.App) {
	appPort := os.Getenv("APP_PORT")
	log.Printf("Starting server on :%s", appPort)
	log.Fatal(app.Listen(fmt.Sprintf(":%s", appPort)))
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
	log.Printf("Stored URL: %s -> %s", shortURL, originalURL)
	response := fmt.Sprintf(`<p>Shortened URL: <a href="/r/%s">%s</a></p>`, shortURL, shortURL)
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

func generateShortURL() string {
	rand.Seed(uint64(time.Now().UnixNano()))
	b := make([]byte, 8)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func setURL(shortURL, originalURL string) error {
	return rdb.Set(ctx, shortURL, originalURL, 0).Err()
}

func getURL(shortURL string) (string, error) {
	return rdb.Get(ctx, shortURL).Result()
}
