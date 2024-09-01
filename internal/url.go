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

func generateShortURL() string {
    rand.Seed(uint64(time.Now().UnixNano()))
    b := make([]byte, 8)
    for i := range b {
        b[i] = letterBytes[rand.Intn(len(letterBytes))]
    }
    return string(b)
}
