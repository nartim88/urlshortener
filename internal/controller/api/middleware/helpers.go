package middleware

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/nartim88/urlshortener/internal/models"
	"github.com/nartim88/urlshortener/pkg/logger"
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
	w          http.ResponseWriter
	zw         *gzip.Writer
	statusCode int
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
	c.statusCode = statusCode
	c.w.Header().Set("Content-Encoding", "gzip")
	c.w.WriteHeader(statusCode)
}

func (c *compressWriter) Close() error {
	if c.statusCode != http.StatusNoContent {
		return c.zw.Close()
	}
	return nil
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

func (c *compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

// decompress распаковывает данные в запросе
func decompress(r *http.Request) error {
	contentEncoding := r.Header.Get("Content-Encoding")
	sendsGzip := strings.Contains(contentEncoding, "gzip")

	if sendsGzip {
		cr, err := newCompressReader(r.Body)
		if err != nil {
			return err
		}

		r.Body = cr
		if err = cr.Close(); err != nil {
			return err
		}
	}
	return nil
}

// canCompress проверяет возможно ли сжимать содержимое ответа
func canCompress(w http.ResponseWriter, r http.Request) bool {
	acceptEncoding := r.Header.Get("Accept-Encoding")
	supportsGzip := strings.Contains(acceptEncoding, "gzip")

	return supportsGzip
}

// setCookieWithToken создает куку с токеном
func setCookieWithToken(rw *http.ResponseWriter, key string) *http.Cookie {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		logger.Log.Error().Stack().Err(err).Msg("error while trying to form uuid")
		http.Error(*rw, err.Error(), http.StatusInternalServerError)
		return nil
	}
	claim := models.Claims{
		RegisteredClaims: jwt.RegisteredClaims{},
		UserID:           newUUID.String(),
	}
	tokenString, err := buildJWTString(claim, key)
	if err != nil {
		logger.Log.Error().Err(err).Msg("error while getting jwt string")
		http.Error(*rw, err.Error(), http.StatusInternalServerError)
		return nil
	}
	cookieName := "token"
	cookie := newCookie(cookieName, tokenString)
	http.SetCookie(*rw, cookie)
	return cookie
}

// validateCookieWithToken проверяет есть ли кука с токеном и валидность
// и возвращает NoCookieWithTokenErr если нет.
func validateCookieWithToken(r http.Request) (*http.Cookie, error) {
	var cookieName = "token"
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return nil, NewNoCookieWithTokenErr(err)
	}
	err = cookie.Valid()
	if err != nil {
		return nil, NewNoCookieWithTokenErr(err)
	}
	return cookie, nil
}

// newCookie возвращает новую куку с предустановленными параметрами
func newCookie(name string, value string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Secure:   true,
		HttpOnly: true,
		//SameSite: http.SameSite(3),
		Path: "/",
	}
}

// buildJWTString создаёт JWT токен и возвращает его в виде строки
func buildJWTString(claims jwt.Claims, key string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// getUserID возвращает ID пользователя из строки с токеном
func getUserID(tokenString string, key string, claims *models.Claims) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(key), nil
		})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf("token is not valid")
	}

	if claims.UserID == "" {
		return "", errors.New("user id is absent in the jwt")
	}
	return claims.UserID, nil
}
