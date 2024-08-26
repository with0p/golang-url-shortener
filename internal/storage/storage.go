package storage

import (
	"context"

	commontypes "github.com/with0p/golang-url-shortener.git/internal/common-types"
)

type Storage interface {
	Read(ctx context.Context, shortURLKey string) (string, error)
	ReadUserID(ctx context.Context, shortURLKey string) (string, error)
	Write(ctx context.Context, userID string, shortURLKey string, fullURL string) error
	WriteBatch(ctx context.Context, userID string, records []commontypes.BatchRecord) error
	SelectAllUserRecords(ctx context.Context, userID string) ([]commontypes.UserRecordData, error)
	MarkAsDeleted(ctx context.Context, shortURLKeys []string) error
}
