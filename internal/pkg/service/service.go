package service

import (
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateRandChars возвращает строку из случайных символов.
func GenerateRandChars(n int) []byte {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	shortKey := make([]byte, n)
	for i := 0; i < n; i++ {
		shortKey[i] = charset[rnd.Intn(len(charset))]
	}
	return shortKey
}
