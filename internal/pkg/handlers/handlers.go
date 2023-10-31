package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nartim88/urlshortener/internal/app/shortener"
	"github.com/nartim88/urlshortener/internal/pkg/logger"
	"github.com/nartim88/urlshortener/internal/pkg/models"
)

// IndexHandle возвращает короткий УРЛ
func IndexHandle(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		logger.Log.Info().Stack().Err(err).Send()
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	body, _ := io.ReadAll(r.Body)

	fURL := models.FullURL(body)
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
	sURL := models.ShortURL(id)
	fURL := shortener.App.Store.Get(sURL)

	if fURL == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Location", string(fURL))
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func JSONGetShortURLHandle(w http.ResponseWriter, r *http.Request) {
	var req models.Request
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)

	if err != nil {
		logger.Log.Info().Err(err).Send()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(buf.Bytes(), &req); err != nil {
		logger.Log.Info().Err(err).Send()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sURL, err := shortener.App.Store.Set(req.FullURL)

	if err != nil {
		logger.Log.Info().Err(err).Send()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res := shortener.App.Configs.BaseURL + "/" + string(sURL)
	resp := models.Response{
		Response: models.ResponsePayload{
			Result: res,
		},
	}

	respDecoded, err := json.Marshal(resp.Response)

	if err != nil {
		logger.Log.Info().Err(err).Send()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(respDecoded)
}
