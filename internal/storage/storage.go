package storage

import (
	"context"

	commontypes "github.com/with0p/golang-url-shortener.git/internal/common-types"
)

type Storage interface {
	Read(ctx context.Context, shortURLKey string) (string, error)
	Write(ctx context.Context, shortURLKey string, fullURL string) error
	WriteBatch(ctx context.Context, records []commontypes.BatchRecord) error
}
