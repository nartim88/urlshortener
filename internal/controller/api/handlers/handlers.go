package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"

	"github.com/nartim88/urlshortener/internal/models"
	"github.com/nartim88/urlshortener/internal/models/api/v1"
	"github.com/nartim88/urlshortener/internal/models/api/v2"
	"github.com/nartim88/urlshortener/internal/service"
	"github.com/nartim88/urlshortener/internal/storage"
	"github.com/nartim88/urlshortener/internal/utils"
	"github.com/nartim88/urlshortener/pkg/logger"
)

const (
	contentType     = "Content-Type"
	textPlain       = "text/plain"
	applicationJSON = "application/json"
)

// IndexHandle возвращает короткий УРЛ
func IndexHandle(svc service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		var sCode = http.StatusCreated
		shortURL, err := svc.CreateShortenURL(ctx, fURL)
		if err != nil {
			var existsErr storage.ErrURLExists
			if errors.As(err, &existsErr) {
				logger.Log.Info().Msgf("%v", existsErr)
				sCode = http.StatusConflict
			} else {
				logger.Log.Info().Err(err).Send()
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		w.Header().Set(contentType, textPlain)
		w.WriteHeader(sCode)
		_, err = w.Write([]byte(*shortURL))
		if err != nil {
			logger.Log.Info().Err(err).Send()
		}
	}
}

// GetURLHandle возвращает полный УРЛ по короткому
func GetURLHandle(svc service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		sID := models.ShortenID(id)

		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		fURL, err := svc.GetFullURL(ctx, sID)
		if err != nil {
			if errors.Is(err, storage.ErrURLDeleted) {
				http.Error(w, err.Error(), http.StatusGone)
				return
			}
			logger.Log.Info().Err(err).Send()
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
}

func GetShortURLHandle(svc service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req v1.Request
		var buf bytes.Buffer

		_, err := buf.ReadFrom(r.Body)

		if err != nil {
			logger.Log.Info().Err(err).Send()
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err = json.Unmarshal(buf.Bytes(), &req); err != nil {
			logger.Log.Info().Err(err).Send()
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		logger.Log.Info().Str("original_url", string(req.FullURL)).Msg("incoming request data:")

		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		var sCode = http.StatusCreated
		shortURL, err := svc.CreateShortenURL(ctx, req.FullURL)
		if err != nil {
			var existsErr storage.ErrURLExists
			if errors.As(err, &existsErr) {
				logger.Log.Info().Msgf("%v", err)
				sCode = http.StatusConflict
			} else {
				logger.Log.Info().Err(err).Send()
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		resp := v1.Response{
			Response: v1.ResponsePayload{
				Result: string(*shortURL),
			},
		}

		respDecoded, err := json.Marshal(resp.Response)
		if err != nil {
			logger.Log.Info().Err(err).Send()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set(contentType, applicationJSON)
		w.WriteHeader(sCode)
		_, err = w.Write(respDecoded)
		if err != nil {
			logger.Log.Info().Err(err).Send()
		}
	}
}

func DBPingHandle(svc service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dsn := svc.GetConfigs().DatabaseDSN
		logger.Log.Info().Str("DATABASE_DSN", dsn).Msg("trying to ping DB via")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		conn, err := pgx.Connect(ctx, dsn)
		if err != nil {
			logger.Log.Error().Stack().Err(err).Send()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer func() {
			if err := conn.Close(ctx); err != nil {
				logger.Log.Error().Stack().Err(err).Send()
			}
		}()

		logger.Log.Info().Msg("ping successfully processed")
		w.WriteHeader(http.StatusOK)
	}
}

// GetBatchShortURLsHandle сохраняет батч урлов
func GetBatchShortURLsHandle(svc service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		var sCode = http.StatusCreated

		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		for _, rData := range req.Data {
			shortURL, err := svc.CreateShortenURL(ctx, rData.FullURL)
			if err != nil {
				var existsErr storage.ErrURLExists
				if errors.As(err, &existsErr) {
					logger.Log.Info().Msgf("%v", existsErr)
					sCode = http.StatusConflict
				} else {
					logger.Log.Info().Err(err).Send()
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
			respPayload = append(respPayload, v2.ResponsePayload{
				CorrelationID: rData.CorrelationID,
				ShortURL:      string(*shortURL),
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
		w.WriteHeader(sCode)
		_, err = w.Write(respDecoded)
		if err != nil {
			logger.Log.Info().Err(err).Msg("error while sending response")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func UserURLsGet(svc service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := utils.GetUserIDFromCtx(r.Context())
		if err != nil {
			logger.Log.Error().Err(err).Send()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		user := models.User{UserID: userID}
		resp, err := svc.GetAllURLs(ctx, user)
		logger.Log.Info().Any("user urls:", resp).Send()
		if err != nil {
			logger.Log.Error().Stack().Err(err).Send()
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(resp) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		respEncoded, err := json.Marshal(resp)
		if err != nil {
			logger.Log.Error().Err(err).Msg("error while serializing response")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set(contentType, applicationJSON)
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(respEncoded)
		if err != nil {
			logger.Log.Info().Err(err).Msg("error while sending response")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	}
}

func UserURLsDelete(svc service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type request struct {
			IDs []models.ShortenID
		}
		var req request

		if err := json.NewDecoder(r.Body).Decode(&req.IDs); err != nil {
			logger.Log.Info().Err(err).Msg("error while decoding request body")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := svc.DeleteURLs(ctx, req.IDs); err != nil {
			logger.Log.Info().Err(err).Msg("error while deleting urls")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}
