package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/nartim88/urlshortener/config"
	"github.com/nartim88/urlshortener/internal/models"
	"github.com/nartim88/urlshortener/pkg/logger"
)

func WithLogging(next http.Handler) http.Handler {
	f := func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()

		uri := r.RequestURI
		method := r.Method

		respData := &responseData{
			status: 0,
			size:   0,
		}
		lrw := loggingResponseWriter{
			ResponseWriter: rw,
			responseData:   respData,
		}

		next.ServeHTTP(&lrw, r)

		logger.Log.Info().
			Str("uri", uri).
			Str("method", method).
			TimeDiff("duration", time.Now(), start).
			Int("status", respData.status).
			Int("size", respData.size).
			Send()
	}
	return http.HandlerFunc(f)
}

func GZipMiddleware(next http.Handler) http.Handler {
	f := func(rw http.ResponseWriter, r *http.Request) {
		canCompr := canCompress(rw, *r)
		logger.Log.Info().Msgf("can compress: %v", canCompr)

		if canCompr {
			cw := newCompressWriter(rw)
			rw = cw
			defer func() {
				if err := cw.Close(); err != nil {
					logger.Log.Info().Err(err).Send()
				}
			}()
		}

		if err := decompress(r); err != nil {
			logger.Log.Info().Err(err).Send()
		}

		next.ServeHTTP(rw, r)
	}
	return http.HandlerFunc(f)
}

func AuthMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		f := func(rw http.ResponseWriter, r *http.Request) {
			cookie, err := validateCookieWithToken(*r)
			if err != nil {
				logger.Log.Info().Msgf("%v", err)
				setCookieWithToken(&rw, cfg.SecretKey)
			} else {
				claims := &models.Claims{}
				tokenString := cookie.Value
				UserID, err := getUserID(tokenString, cfg.SecretKey, claims)
				if err != nil {
					logger.Log.Error().Err(err).Send()
					http.Error(rw, err.Error(), http.StatusUnauthorized)
					return
				}
				k := models.UserIDCtxKey("userID")
				ctx := context.WithValue(r.Context(), k, UserID)
				r = r.WithContext(ctx)
			}
			next.ServeHTTP(rw, r)
		}
		return http.HandlerFunc(f)
	}
}
