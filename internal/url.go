package internal

import (
	"context"
	"sync"
	"time"

	"golang.org/x/exp/rand"
)

var (
	mu            sync.RWMutex
	memoryStorage map[string]string
	ctx           = context.Background()
)

// generateShortURL generates a random 5-character string using the characters
// in the letterBytes slice. The string is suitable for use as a short URL.
func generateShortURL() string {
	rand.Seed(uint64(time.Now().UnixNano()))
	b := make([]byte, 5)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
