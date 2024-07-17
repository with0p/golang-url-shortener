package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

var logger = zap.Must(zap.NewProduction()).Sugar()

func HandleWithLogging(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		extendedResponseData := &extendedResponseData{}

		extendedW := &extendedResponseWriter{
			ResponseWriter:       w,
			extendedResponseData: extendedResponseData,
		}

		handler.ServeHTTP(extendedW, r)

		duration := time.Since(start)

		logger.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"duration", duration,
			"status", extendedResponseData.statusCode,
			"size", extendedResponseData.size,
		)
	}
}
