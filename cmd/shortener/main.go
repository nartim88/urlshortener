package main

import (
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var URLs = make(map[string]string)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const host = "localhost"
const port = "8080"
const schema = "http"
const shortURLLen = 8

type fullURL struct {
	name string
}
type shortURL string

type Repository interface {
	// Save сохранение данных.
	Save(storage any) shortURL

	// IsExist проверка есть ли короткий урл в базе.
	// True если есть, false - нет.
	IsExist(shortURL string, storage any) bool
}

func (url *fullURL) Save(storage any) shortURL {
	sURL := generateShortURL(shortURLLen)
	s := storage.(map[string]string)
	s[sURL] = url.name
	return shortURL(sURL)
}

func (url *fullURL) IsExist(shortURL string, storage any) bool {
	s := storage.(map[string]string)
	if _, ok := s[shortURL]; !ok {
		return false
	}
	return true
}

func GetFullURLbyShortURL(shortURL string, storage any) string {
	s := storage.(map[string]string)
	val, ok := s[shortURL]
	if !ok {
		return ""
	}
	return val
}

// generateShortURL возвращает строку из случайных символов.
func generateShortURL(n int64) string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	shortKey := make([]byte, n)
	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortKey)
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	body, _ := io.ReadAll(r.Body)

	switch r.Method {
	case http.MethodPost:
		fURL := fullURL{string(body)}
		sURL := fURL.Save(URLs)
		result := schema + "://" + host + ":" + port + "/" + sURL
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(result))
	case http.MethodGet:
		path := r.URL.Path
		id := strings.Split(path, "/")[1]
		fURL := GetFullURLbyShortURL(id, URLs)
		w.Header().Set("Location", fURL)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusTemporaryRedirect)
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", mainPage)

	err := http.ListenAndServe(host+":"+port, mux)
	if err != nil {
		log.Fatal(err)
	}
}
