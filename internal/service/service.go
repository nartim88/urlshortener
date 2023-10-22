package service

import (
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateRandChars возвращает строку из случайных символов.
func GenerateRandChars(n int64) []byte {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	shortKey := make([]byte, n)
	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}
	return shortKey
}
