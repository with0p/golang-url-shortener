package logger

import "net/http"

type extendedResponseData struct {
	statusCode int
	size       int
}

type extendedResponseWriter struct {
	http.ResponseWriter
	extendedResponseData *extendedResponseData
}

func (r *extendedResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.extendedResponseData.size = size
	return size, err
}

func (r *extendedResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.extendedResponseData.statusCode = statusCode
}
