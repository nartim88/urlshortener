package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nartim88/urlshortener/internal/app/shortener"
	"github.com/nartim88/urlshortener/internal/pkg/routers"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
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
	shortener.New() // init variables for test

	ts := httptest.NewServer(routers.MainRouter())
	defer ts.Close()

	type want struct {
		contentType string
		statusCode  int
	}

	var testTable = []struct {
		url    string
		method string
		want   want
	}{
		{
			url:    "/",
			method: http.MethodPost,
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusCreated,
			},
		},
		{
			url:    "/HMOUQTFX",
			method: http.MethodGet,
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
	}

	for _, v := range testTable {
		resp, _ := testRequest(t, ts, v.method, v.url)
		_ = resp.Body.Close() // для прохождения теста

		assert.Equal(t, v.want.statusCode, resp.StatusCode)
		assert.Equal(t, v.want.contentType, resp.Header.Get("Content-Type"))
	}
}

func TestAPI(t *testing.T) {
	srv := httptest.NewServer(routers.MainRouter())
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
