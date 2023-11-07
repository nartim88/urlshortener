package main

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nartim88/urlshortener/internal/app/shortener"
	"github.com/nartim88/urlshortener/internal/pkg/handlers"
	"github.com/nartim88/urlshortener/internal/pkg/routers"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)

	defer func() {
		_ = resp.Body.Close()
	}()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestMainRouter(t *testing.T) {
	shortener.App.Init() // init variables for test

	ts := httptest.NewServer(routers.MainRouter())
	defer ts.Close()

	type want struct {
		contentType string
		statusCode  int
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
				statusCode:  http.StatusCreated,
			},
		},
		{
			name:   "url_not_found",
			url:    "/HMOUQTFX",
			method: http.MethodGet,
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buf := bytes.NewBufferString(tc.body)
			resp, _ := testRequest(t, ts, tc.method, tc.url, buf)
			defer resp.Body.Close()

			assert.Equal(t, tc.want.statusCode, resp.StatusCode)
			assert.Equal(t, tc.want.contentType, resp.Header.Get("Content-Type"))
		})
	}
}

func TestAPI(t *testing.T) {
	handler := http.HandlerFunc(handlers.JSONGetShortURLHandle)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	type want struct {
		contentType string
		statusCode  int
	}

	var testCases = []struct {
		name   string
		path   string
		method string
		body   string
		want   want
	}{
		{
			name:   "method_post",
			path:   "/api/shorten",
			method: http.MethodPost,
			body:   `{"url": "https://ya.ru"}`,
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusCreated,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := resty.New().R()

			req.Method = tc.method
			req.URL = srv.URL + tc.path
			req.Body = tc.body
			req.Header.Set("Content-Type", "application/json")

			resp, err := req.Send()

			assert.NoError(t, err)
			assert.Equal(t, tc.want.statusCode, resp.StatusCode())
			assert.Equal(t, tc.want.contentType, resp.Header().Get("Content-Type"))
		})
	}
}

func TestGzipCompression(t *testing.T) {
	srv := httptest.NewServer(routers.MainRouter())
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

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp.StatusCode)

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

		resp, err := req.Send()

		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp.StatusCode())

		zr, err := gzip.NewReader(resp.RawBody())
		defer resp.RawBody().Close()
		require.NoError(t, err)

		_, err = io.ReadAll(zr)
		require.NoError(t, err)
	})
}
