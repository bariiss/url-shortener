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
    rdb           *redis.Client
    ctx           = context.Background()
    mu            sync.RWMutex
    redisActive   bool
    memoryStorage map[string]string
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func loadEnv() {
    if err := godotenv.Load(); err != nil {
        log.Fatalf("Error loading .env file")
    }
    redisActive = os.Getenv("REDIS_ACTIVE") == "true"
}

func initRedis() {
    if redisActive {
        redisAddr := os.Getenv("REDIS_ADDR")
        rdb = redis.NewClient(&redis.Options{
            Addr: redisAddr,
        })

        if _, err := rdb.Ping(ctx).Result(); err != nil {
            log.Fatalf("Error connecting to Redis: %v", err)
        }
        log.Println("Connected to Redis")
    } else {
        memoryStorage = make(map[string]string)
        log.Println("Redis is not active. Using in-memory storage.")
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

func showWelcomeMessage() {
    fmt.Println(`
    ====================================
    Welcome to the URL Shortener Service
    ====================================
    `)
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

func generateShortURL() string {
    rand.Seed(uint64(time.Now().UnixNano()))
    b := make([]byte, 8)
    for i := range b {
        b[i] = letterBytes[rand.Intn(len(letterBytes))]
    }
    return string(b)
}

func setURL(shortURL, originalURL string) error {
    if redisActive {
        return rdb.Set(ctx, shortURL, originalURL, 0).Err()
    } else {
        mu.Lock()
        defer mu.Unlock()
        memoryStorage[shortURL] = originalURL
        return nil
    }
}

func getURL(shortURL string) (string, error) {
    if redisActive {
        return rdb.Get(ctx, shortURL).Result()
    } else {
        mu.RLock()
        defer mu.RUnlock()
        originalURL, exists := memoryStorage[shortURL]
        if !exists {
            return "", fmt.Errorf("URL not found")
        }
        return originalURL, nil
    }
}
