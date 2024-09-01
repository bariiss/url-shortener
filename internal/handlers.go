package internal

import (
    "fmt"
    "log"

    "github.com/gofiber/fiber/v2"
)

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
