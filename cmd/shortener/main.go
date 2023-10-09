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
	IsExist(shortUrl string) bool
}

func (url *fullURL) Save(storage any) shortURL {
	sUrl := generateShortUrl(shortURLLen)
	s := storage.(map[string]string)
	s[sUrl] = url.name
	return shortURL(sUrl)
}

func (url *fullURL) IsExist(shortUrl string) bool {
	if _, ok := URLs[shortUrl]; !ok {
		return false
	}
	return true
}

// generateShortUrl возвращает строку из случайных символов.
func generateShortUrl(n int64) string {
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

	if r.Method == http.MethodPost {
		fURL := fullURL{string(body)}
		sURL := fURL.Save(URLs)
		result := schema + "://" + host + ":" + port + "/" + sURL

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)

		_, _ = w.Write([]byte(result))

	} else if r.Method == http.MethodGet {
		path := r.URL.Path
		id := strings.Split(path, "/")[1]

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusTemporaryRedirect)

		_, _ = w.Write([]byte(id))

	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", mainPage)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
