package compressor

import (
	"compress/gzip"
	"net/http"
)

type CompressorWriter struct {
	httpWriter http.ResponseWriter
	gzipWriter *gzip.Writer
}

func newCompressorWriter(httpWriter http.ResponseWriter) *CompressorWriter {
	return &CompressorWriter{
		httpWriter: httpWriter,
		gzipWriter: gzip.NewWriter(httpWriter),
	}
}

func (compressorWriter *CompressorWriter) Header() http.Header {
	return compressorWriter.httpWriter.Header()
}

func (compressorWriter *CompressorWriter) WriteHeader(statusCode int) {
	compressorWriter.httpWriter.Header().Set("Content-Encoding", "gzip")
	compressorWriter.httpWriter.WriteHeader(statusCode)
}

func (compressorWriter CompressorWriter) Write(data []byte) (int, error) {
	return compressorWriter.gzipWriter.Write(data)
}

func (compressorWriter CompressorWriter) Close() error {
	return compressorWriter.gzipWriter.Close()
}
