package main

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nartim88/urlshortener/config"
	"github.com/nartim88/urlshortener/internal/app/shortener"
	"github.com/nartim88/urlshortener/internal/controller/api/middleware"
	"github.com/nartim88/urlshortener/internal/controller/api/routers"
)

func initApp() shortener.Application {
	cfg := &config.Config{
		SecretKey: "secret_key",
	}
	app := shortener.NewApp()
	app.Init(cfg) // init variables for test
	return *app
}

// newCookie возвращает новую куку с предустановленными параметрами
func newCookie(name string, value string) *http.Cookie {
	return &http.Cookie{
		Name:  name,
		Value: value,
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

// getCookieWithToken создает куку с токеном
func getCookieWithToken(key string) (*http.Cookie, error) {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	claim := middleware.Claims{
		RegisteredClaims: jwt.RegisteredClaims{},
		UserID:           newUUID.String(),
	}
	tokenString, err := buildJWTString(claim, key)
	if err != nil {
		return nil, err
	}
	cookieName := "token"
	cookie := newCookie(cookieName, tokenString)
	return cookie, nil
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader, key string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)

	cookie, err := getCookieWithToken(key)
	require.NoError(t, err)

	req.AddCookie(cookie)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)

	defer func() {
		_ = resp.Body.Close()
	}()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

// testJSONRequest отправляет json request с кукой с токеном
func testJSONRequest(t *testing.T, method, url, path string, body any, key string) *resty.Response {
	req := resty.New().R()
	req.Method = method
	req.URL = url + path
	req.Body = body
	req.Header.Set("Content-Type", "application/json")

	cookie, err := getCookieWithToken(key)
	require.NoError(t, err)

	req.SetCookie(cookie)

	resp, err := req.Send()
	require.NoError(t, err)

	return resp
}

func TestMainRouter(t *testing.T) {
	app := initApp()
	ts := httptest.NewServer(routers.MainRouter(app.Service))
	defer ts.Close()

	type want struct {
		contentType string
		statusCodes []int
	}

	var testCases = []struct {
		name   string
		url    string
		method string
		body   string
		want   want
	}{
		{
			name:   "success_request",
			url:    "/",
			method: http.MethodPost,
			body:   "https://ya.ru",
			want: want{
				contentType: "text/plain",
				statusCodes: []int{201, 409},
			},
		},
		{
			name:   "url_not_found",
			url:    "/HMOUQTFX",
			method: http.MethodGet,
			want: want{
				statusCodes: []int{404},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buf := bytes.NewBufferString(tc.body)
			resp, _ := testRequest(t, ts, tc.method, tc.url, buf, app.Configs.SecretKey)
			defer resp.Body.Close()

			assert.Contains(t, tc.want.statusCodes, resp.StatusCode)
			assert.Equal(t, tc.want.contentType, resp.Header.Get("Content-Type"))
		})
	}

	t.Run("/{id}_GET_temp_redirect", func(t *testing.T) {
		// создаем новый короткий урл
		reqText := "https://ya.ru"
		buf := bytes.NewBufferString(reqText)
		_, body := testRequest(t, ts, http.MethodPost, "/", buf, app.Configs.SecretKey)

		// проверяем вторым запросом, что он создан и приходит правильный ответ
		resp := testJSONRequest(t, http.MethodGet, ts.URL, body, nil, app.Configs.SecretKey)
		assert.Equal(t, http.StatusOK, resp.StatusCode())
	})
}

func TestAPI(t *testing.T) {
	app := initApp()
	srv := httptest.NewServer(routers.MainRouter(app.Service))
	defer srv.Close()

	type want struct {
		contentType string
		statusCodes []int
	}

	var testCases = []struct {
		name   string
		path   string
		method string
		body   string
		want   want
	}{
		{
			name:   "/api/shorten_POST",
			path:   "/api/shorten",
			method: http.MethodPost,
			body:   `{"url": "https://ya.ru"}`,
			want: want{
				contentType: "application/json",
				statusCodes: []int{201, 409},
			},
		},
		{
			name:   "/api/shorten/batch_POST",
			path:   "/api/shorten/batch",
			method: http.MethodPost,
			body: `
				[
					{
						"correlation_id": "salt1",
						"original_url": "https://ya.ru"
					},
					{
						"correlation_id": "salt2",
						"original_url": "https://google.ru"
					}
				]`,
			want: want{
				contentType: "application/json",
				statusCodes: []int{201, 409},
			},
		},
		{
			name:   "/api/user/urls_GET",
			path:   "/api/user/urls",
			method: http.MethodGet,
			want: want{
				statusCodes: []int{200, 204},
			},
		},
		{
			name:   "/api/user/urls_DELETE",
			path:   "/api/user/urls",
			method: http.MethodDelete,
			body:   `["xLbDjDFR", "SnDb3L2k"]`,
			want: want{
				statusCodes: []int{202},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp := testJSONRequest(t, tc.method, srv.URL, tc.path, tc.body, app.Configs.SecretKey)

			assert.Contains(t, tc.want.statusCodes, resp.StatusCode())
			assert.Equal(t, tc.want.contentType, resp.Header().Get("Content-Type"))

			t.Logf("response text: %s", resp.Body())
		})
	}
}

func TestGzipCompression(t *testing.T) {
	app := initApp()
	srv := httptest.NewServer(routers.MainRouter(app.Service))
	defer srv.Close()

	requestBody := `{"url": "https://ya.ru"}`

	t.Run("sends_gzip", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		zb := gzip.NewWriter(buf)
		_, err := zb.Write([]byte(requestBody))
		require.NoError(t, err)
		err = zb.Close()
		require.NoError(t, err)
		target := srv.URL + "/api/shorten"

		req := httptest.NewRequest(http.MethodPost, target, buf)
		req.RequestURI = ""
		req.Header.Set("Content-Encoding", "gzip")
		req.Header.Set("Content-Type", "application/json")

		cookie, err := getCookieWithToken(app.Configs.SecretKey)
		require.NoError(t, err)

		req.AddCookie(cookie)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		require.Contains(t, []int{201, 409}, resp.StatusCode)

		defer resp.Body.Close()

		_, err = io.ReadAll(resp.Body)
		require.NoError(t, err)
	})

	t.Run("accepts_gzip", func(t *testing.T) {
		req := resty.New().R()
		req.URL = srv.URL + "/api/shorten"
		req.Method = http.MethodPost
		req.Body = requestBody
		req.Header.Set("Accept-Encoding", "gzip")
		req.Header.Set("Content-Type", "application/json")
		req.SetDoNotParseResponse(true)

		cookie, err := getCookieWithToken(app.Configs.SecretKey)
		require.NoError(t, err)

		req.SetCookie(cookie)

		resp, err := req.Send()
		require.NoError(t, err)
		require.Contains(t, []int{201, 409}, resp.StatusCode())

		zr, err := gzip.NewReader(resp.RawBody())
		defer resp.RawBody().Close()
		require.NoError(t, err)

		_, err = io.ReadAll(zr)
		require.NoError(t, err)
	})
}
