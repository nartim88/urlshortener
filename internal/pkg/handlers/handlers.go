package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/nartim88/urlshortener/internal/app/shortener"
	"github.com/nartim88/urlshortener/internal/pkg/logger"
	"github.com/nartim88/urlshortener/internal/pkg/models"
	"github.com/nartim88/urlshortener/internal/pkg/models/api/v1"
	v2 "github.com/nartim88/urlshortener/internal/pkg/models/api/v2"
)

const (
	contentType     = "Content-Type"
	textPlain       = "text/plain"
	applicationJSON = "application/json"
)

// IndexHandle возвращает короткий УРЛ
func IndexHandle(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		logger.Log.Info().Stack().Err(err).Send()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Log.Info().Err(err).Send()
	}

	if len(body) == 0 {
		err = errors.New("request Body is empty")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fURL := models.FullURL(body)
	sID, err := shortener.App.Store.Set(fURL)

	if err != nil {
		logger.Log.Info().Stack().Err(err).Send()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortURL := shortener.App.Configs.BaseURL + "/" + string(*sID)

	w.Header().Set(contentType, textPlain)
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(shortURL))

	if err != nil {
		logger.Log.Info().Err(err).Send()
	}
}

// GetURLHandle возвращает полный УРЛ по короткому
func GetURLHandle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	sID := models.ShortenID(id)
	fURL, err := shortener.App.Store.Get(sID)

	if err != nil {
		logger.Log.Info().Err(err).Send()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if fURL == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Location", string(*fURL))
	w.Header().Set(contentType, textPlain)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func GetShortURLHandle(w http.ResponseWriter, r *http.Request) {
	var req v1.Request
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
	logger.Log.Info().Str("original_url", string(req.FullURL)).Msg("incoming request data:")

	sID, err := shortener.App.Store.Set(req.FullURL)
	if err != nil {
		logger.Log.Info().Err(err).Send()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortURL := shortener.App.Configs.BaseURL + "/" + string(*sID)
	resp := v1.Response{
		Response: v1.ResponsePayload{
			Result: shortURL,
		},
	}

	respDecoded, err := json.Marshal(resp.Response)
	if err != nil {
		logger.Log.Info().Err(err).Send()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set(contentType, applicationJSON)
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(respDecoded)
	if err != nil {
		logger.Log.Info().Err(err).Send()
	}
}

func DBPingHandle(w http.ResponseWriter, r *http.Request) {
	logger.Log.Info().Str("DATABASE_DSN", shortener.App.Configs.DatabaseDSN).Msg("trying to ping DB via")

	conn, err := pgx.Connect(context.Background(), shortener.App.Configs.DatabaseDSN)
	if err != nil {
		logger.Log.Error().Stack().Err(err).Send()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := conn.Close(context.Background()); err != nil {
			logger.Log.Error().Stack().Err(err).Send()
		}
	}()

	logger.Log.Info().Msg("ping successfully processed")
	w.WriteHeader(http.StatusOK)
}

func GetBatchShortURLsHandle(w http.ResponseWriter, r *http.Request) {
	var req v2.Request
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		logger.Log.Error().Err(err).Msg("error while reading from request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &req.Data); err != nil {
		logger.Log.Error().Err(err).Msg("error while deserializing json")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger.Log.Info().Msgf("got batch: %+v", req)

	var respPayload []v2.ResponsePayload

	for _, rData := range req.Data {
		sID, err := shortener.App.Store.Set(rData.FullURL)
		if err != nil {
			logger.Log.Error().Err(err).Msg("error while saving data to db")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		shortURL := shortener.App.Configs.BaseURL + "/" + string(*sID)
		respPayload = append(respPayload, v2.ResponsePayload{
			CorrelationID: rData.CorrelationID,
			ShortURL:      shortURL,
		})
	}

	resp := v2.Response{Response: respPayload}

	respDecoded, err := json.Marshal(resp.Response)
	if err != nil {
		logger.Log.Error().Err(err).Msg("error while serializing response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set(contentType, applicationJSON)
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(respDecoded)
	if err != nil {
		logger.Log.Info().Err(err).Msg("error while sending response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
