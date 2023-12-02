package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/nartim88/urlshortener/internal/app/shortener"
	"github.com/nartim88/urlshortener/internal/pkg/config"
	"github.com/nartim88/urlshortener/internal/pkg/logger"
	"github.com/nartim88/urlshortener/internal/pkg/service"
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
		if canCompress(*r) {
			logger.Log.Info().Msgf("can compress %v", canCompress(*r))
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

func AuthMiddleware(next http.Handler) http.Handler {
	f := func(rw http.ResponseWriter, r *http.Request) {
		cookie, err := validateCookieWithToken(*r)
		if err != nil {
			logger.Log.Info().Msgf("%v", err)
			setCookieWithToken(&rw)
		} else {
			claims := &config.Claims{}
			tokenString := cookie.Value
			UserID, err := service.GetUserId(tokenString, shortener.App.Configs.SecretKey, claims)
			if err != nil {
				logger.Log.Error().Err(err).Send()
				http.Error(rw, err.Error(), http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), "userID", UserID)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(rw, r)
	}
	return http.HandlerFunc(f)
}
