package localfile

import "github.com/google/uuid"

type LocalFileRecord struct {
	UUID        string `json:"uuid"`
	UserId      string `json:"user_id"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func NewLocalFileRecord(userId string, key string, value string) *LocalFileRecord {
	return &LocalFileRecord{
		UUID:        uuid.New().String(),
		UserId:      userId,
		ShortURL:    key,
		OriginalURL: value,
	}
}
