package handlers

import (
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nartim88/urlshortener/internal/pkg/logger"
	"github.com/nartim88/urlshortener/internal/pkg/storage"

	"github.com/nartim88/urlshortener/internal/app/shortener"
)

// IndexHandle возвращает короткий УРЛ
func IndexHandle(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		logger.Log.Info().Stack().Err(err).Send()
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	body, _ := io.ReadAll(r.Body)

	fURL := storage.FullURL(body)
	sURL, err := shortener.App.Store.Set(fURL)
	if err != nil {
		logger.Log.Info().Stack().Err(err).Send()
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	result := shortener.App.Configs.BaseURL + "/" + string(sURL)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(result))
}

// GetURLHandle возвращает полный УРЛ по короткому
func GetURLHandle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	sURL := storage.ShortURL(id)
	fURL := shortener.App.Store.Get(sURL)
	if fURL == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Location", string(fURL))
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
