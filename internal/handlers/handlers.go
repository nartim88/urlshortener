package handlers

import (
	"github.com/nartim88/urlshortener/internal/app/shortener"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/nartim88/urlshortener/internal/storage"
)

var URLs = make(map[storage.ShortURL]storage.FullURL)

// IndexHandle возвращает короткий УРЛ
func IndexHandle(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	body, _ := io.ReadAll(r.Body)

	fURL := storage.FullURL(body)
	sURL := fURL.Save(URLs)
	result := shortener.CFG.BaseURL + "/" + string(sURL)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(result))
}

// GetURLHandle возвращает полный УРЛ по короткому
func GetURLHandle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	sURL := storage.ShortURL(id)

	fURL := shortener.Get(&sURL, URLs)
	if fURL == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Location", string(fURL))
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
