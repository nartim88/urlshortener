package middleware

import (
	"net/http"
	"time"

	"github.com/nartim88/urlshortener/internal/pkg/logger"
)

var All = []func(http.Handler) http.Handler{
	WithLogging,
	GZipMiddleware,
}

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
		if CanCompress(*r) {
			cw := newCompressWriter(rw)
			rw = cw
			defer func() {
				if err := cw.Close(); err != nil {
					logger.Log.Info().Err(err).Send()
				}
			}()
		}

		if err := Decompress(r); err != nil {
			logger.Log.Info().Err(err).Send()
		}

		next.ServeHTTP(rw, r)
	}
	return http.HandlerFunc(f)
}
