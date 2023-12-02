package service

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/nartim88/urlshortener/internal/pkg/config"
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

// BuildJWTString создаёт JWT токен и возвращает его в виде строки
func BuildJWTString(claims jwt.Claims, key string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GetUserId возвращает ID пользователя из токена
func GetUserId(tokenString string, key string, claims *config.Claims) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(key), nil
		})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf("token is not valid")
	}

	if claims.UserID == "" {
		return "", errors.New("user id is absent in the jwt")
	}
	return claims.UserID, nil
}

func NewCookie(name string, value string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSite(3),
	}
}
