package middleware

import "net/http"

type (
	responseData struct {
		status int
		size   int
	}
	loggingResponseData struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseData) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseData) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}
