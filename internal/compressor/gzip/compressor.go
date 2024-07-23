package compressor

import (
	"fmt"
	"net/http"
	"strings"
)

func HandleWithGzipCompressor(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writer := w

		fmt.Println("Accept-Encoding", r.Header.Get("Accept-Encoding"))
		fmt.Println("content", r.Header.Get("content-type"))
		acceptsEncoding := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
		validContentType := strings.Contains("text/plain, application/json", r.Header.Get("content-type"))

		if validContentType && acceptsEncoding {
			compressorWriter := newCompressorWriter(w)
			writer = compressorWriter
			defer compressorWriter.Close()
		}

		if strings.Contains(w.Header().Get("Content-Encoding"), "gzip") {
			compressorReader, err := NewCompressorReader(r.Body)

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}

			r.Body = compressorReader
			defer compressorReader.Close()
		}

		handler.ServeHTTP(writer, r)

	}
}
