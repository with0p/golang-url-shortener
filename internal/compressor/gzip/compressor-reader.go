package compressor

import (
	"compress/gzip"
	"io"
)

type CompressorReader struct {
	reader     io.Reader
	gzipReader *gzip.Reader
}

func NewCompressorReader(reader io.ReadCloser) (*CompressorReader, error) {
	gzipReader, err := gzip.NewReader(reader)

	if err != nil {
		return nil, err
	}

	return &CompressorReader{
		reader:     reader,
		gzipReader: gzipReader,
	}, nil
}

func (compressorReader CompressorReader) Read(data []byte) (int, error) {
	return compressorReader.gzipReader.Read(data)
}

func (compressorReader CompressorReader) Close() error {
	return compressorReader.gzipReader.Close()
}
