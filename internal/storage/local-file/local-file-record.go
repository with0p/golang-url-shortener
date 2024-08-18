package localfile

import "github.com/google/uuid"

type LocalFileRecord struct {
	UUID        string `json:"uuid"`
	UserID      string `json:"user_id"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func NewLocalFileRecord(userID string, key string, value string) *LocalFileRecord {
	return &LocalFileRecord{
		UUID:        uuid.New().String(),
		UserID:      userID,
		ShortURL:    key,
		OriginalURL: value,
	}
}
