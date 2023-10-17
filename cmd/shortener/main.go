package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/nartim88/urlshortener/cmd/shortener/config"
)

var URLs = make(map[ShortURL]FullURL)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const shortURLLen = 8

type FullURL string
type ShortURL string

type Saver interface {
	// Save сохранение данных.
	Save(storage any) ShortURL

	// IsExist проверка есть ли короткий урл в базе.
	// True если есть, false - нет.
	IsExist(shortURL string, storage any) bool
}

func (url *FullURL) Save(storage any) ShortURL {
	sURL := GenerateShortURL(shortURLLen)
	s := storage.(map[ShortURL]FullURL)
	s[sURL] = *url
	return sURL
}

func (url *FullURL) IsExist(shortURL string, storage any) bool {
	s := storage.(map[string]string)
	if _, ok := s[shortURL]; !ok {
		return false
	}
	return true
}

// GetFullURLbyShortURL возвращает полный урл по сокращенному.
func GetFullURLbyShortURL(sURL *ShortURL, storage any) FullURL {
	s := storage.(map[ShortURL]FullURL)
	val, ok := s[*sURL]
	if !ok {
		return ""
	}
	return val
}

// GenerateShortURL возвращает строку из случайных символов.
func GenerateShortURL(n int64) ShortURL {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	shortKey := make([]byte, n)
	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}
	return ShortURL(shortKey)
}

// indexHandle возвращает короткий УРЛ
func indexHandle(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	body, _ := io.ReadAll(r.Body)

	fURL := FullURL(body)
	sURL := fURL.Save(URLs)
	result := config.CFG.BaseURL + "/" + string(sURL)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(result))
}

// getURLHandle возвращает полный УРЛ по короткому
func getURLHandle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	sURL := ShortURL(id)

	fURL := GetFullURLbyShortURL(&sURL, URLs)
	if fURL == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Location", string(fURL))
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func mainRouter() chi.Router {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Post("/", indexHandle)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", getURLHandle)
		})
	})

	return r
}

func main() {
	config.InitConfigs()

	fmt.Printf("Runnig server on %s", config.CFG.RunAddr)

	err := http.ListenAndServe(config.CFG.RunAddr, mainRouter())

	if err != nil {
		log.Fatal(err)
	}
}
