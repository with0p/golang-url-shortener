package middlewares

import (
	"net/http"

	"github.com/with0p/golang-url-shortener.git/internal/compressor/gzip"
	"github.com/with0p/golang-url-shortener.git/internal/logger"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func conveyor(h http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}

func UseMiddlewares(handler http.HandlerFunc) http.HandlerFunc {
	return conveyor(handler, compressor.HandleWithGzipCompressor, logger.HandleWithLogging)
}
