package storage

import (
	"context"

	commontypes "github.com/with0p/golang-url-shortener.git/internal/common-types"
)

type Storage interface {
	Read(ctx context.Context, shortURLKey string) (string, error)
	Write(ctx context.Context, userId string, shortURLKey string, fullURL string) error
	WriteBatch(ctx context.Context, userId string, records []commontypes.BatchRecord) error
	SelectAllUserRecords(ctx context.Context, userId string) ([]commontypes.UserRecordData, error)
}
