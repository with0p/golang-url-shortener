package localfile

import "github.com/google/uuid"

type LocalFileRecord struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func NewLocalFileRecord(key string, value string) *LocalFileRecord {
	return &LocalFileRecord{
		UUID:        uuid.New().String(),
		ShortURL:    key,
		OriginalURL: value,
	}
}
