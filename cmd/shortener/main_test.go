package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeRequest(
	method string,
	target string,
	body io.Reader,
	handler func(rw http.ResponseWriter, r *http.Request),
) *http.Response {

	request := httptest.NewRequest(method, target, body)

	w := httptest.NewRecorder()
	h := http.HandlerFunc(handler)
	h(w, request)

	return w.Result()
}

func Test_mainPage(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		location    string
		hostAndPort string
		shortenLen  int
		contentLen  string
	}
	tests := []struct {
		name        string
		method      string
		want        want
		requestBody string
		target      string
	}{
		{
			name:   "mainPage POST",
			target: "/",
			method: http.MethodPost,
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusCreated,
				hostAndPort: "localhost:8080",
				shortenLen:  8,
			},
			requestBody: "ya.ru",
		},
		{
			name:   "mainPage GET",
			target: "/",
			method: http.MethodGet,
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusTemporaryRedirect,
				location:    "https://ya.ru",
				contentLen:  "",
			},
		},
		{
			name:   "mainPage error",
			target: "/",
			method: http.MethodDelete,
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.method {

			case http.MethodPost:
				body := strings.NewReader(tt.requestBody)
				result := makeRequest(tt.method, tt.target, body, mainPage)

				assert.Equal(t, tt.want.statusCode, result.StatusCode)
				assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

				respBody, err := io.ReadAll(result.Body)
				require.NoError(t, err)
				err = result.Body.Close()
				require.NoError(t, err)

				split := strings.Split(string(respBody), "/")
				hostAndPort := split[2]
				id := split[3]

				assert.Equal(t, tt.want.hostAndPort, hostAndPort)
				assert.NotEmpty(t, id)
				assert.Len(t, id, tt.want.shortenLen)

			case http.MethodGet:
				result := makeRequest(tt.method, tt.target, nil, mainPage)

				assert.Equal(t, tt.want.statusCode, result.StatusCode)
				assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
				assert.Equal(t, tt.want.location, result.Header.Get("Location"))
				assert.Equal(t, tt.want.contentLen, result.Header.Get("Content-Length"))

			default:
				result := makeRequest(tt.method, tt.target, nil, mainPage)

				assert.Equal(t, tt.want.statusCode, result.StatusCode)
			}
		})
	}
}
