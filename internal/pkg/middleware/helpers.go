package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/nartim88/urlshortener/internal/pkg/logger"
)

type (
	responseData struct {
		status int
		size   int
	}
	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// compressWriter реализует интерфейс http.ResponseWriter,
// сжимает передаваемые данные и выставляет правильные HTTP-заголовки
type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}
func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

func (c *compressWriter) Close() error {
	return c.zw.Close()
}

// compressReader реализует интерфейс io.ReadCloser и
// декомпрессирует получаемые от клиента данные
type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

func Compress(rw http.ResponseWriter, r *http.Request) http.ResponseWriter {
	contentType := r.Header.Get("Content-Type")
	canCompress := strings.Contains(contentType, "text/html")
	canCompress = strings.Contains(contentType, "application/json")
	acceptEncoding := r.Header.Get("Accept-Encoding")
	supportsGzip := strings.Contains(acceptEncoding, "gzip")

	if canCompress && supportsGzip {
		cw := newCompressWriter(rw)
		if err := cw.Close(); err != nil {
			logger.Log.Info().Err(err).Send()
		}
		return cw
	}
	return rw
}

func Decompress(r *http.Request) error {
	contentEncoding := r.Header.Get("Content-Encoding")
	sendsGzip := strings.Contains(contentEncoding, "gzip")

	if sendsGzip {
		cr, err := newCompressReader(r.Body)
		if err != nil {
			return err
		}

		r.Body = cr
		if err := cr.Close(); err != nil {
			return err
		}
	}
	return nil
}
