package compressor

import (
	"net/http"
	"strings"

	"github.com/with0p/golang-url-shortener.git/internal/logger"
)

func HandleWithGzipCompressor(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writer := w

		encodingAccepted := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
		validContentType := strings.Contains("text/html, application/json", r.Header.Get("content-type"))

		if encodingAccepted && validContentType {
			compressorWriter := newCompressorWriter(w)
			writer = compressorWriter
			defer compressorWriter.Close()
		}

		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			compressorReader, err := newCompressorReader(r.Body)

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				logger.LogError(err)
			}

			r.Body = compressorReader
			defer compressorReader.Close()
		}

		handler.ServeHTTP(writer, r)

	}
}
