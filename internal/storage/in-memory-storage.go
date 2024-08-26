package storage

import (
	"context"
	"errors"

	commontypes "github.com/with0p/golang-url-shortener.git/internal/common-types"
)

type URLStorageMap map[string]string

type InMemoryStorage struct {
	urlMap URLStorageMap
}

func NewInMemoryStorage(storageMap URLStorageMap) *InMemoryStorage {
	return &InMemoryStorage{
		urlMap: storageMap,
	}
}

func (storage *InMemoryStorage) Write(ctx context.Context, userID string, shortURLKey string, fullURL string) error {
	storage.urlMap[shortURLKey] = fullURL
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func (storage *InMemoryStorage) WriteBatch(ctx context.Context, userID string, records []commontypes.BatchRecord) error {
	for _, r := range records {
		storage.urlMap[r.ShortURLKey] = r.FullURL
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func (storage *InMemoryStorage) Read(ctx context.Context, shortURLKey string) (string, error) {
	fullURL, ok := storage.urlMap[shortURLKey]

	if !ok {
		return "", errors.New("not found")
	}

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		return fullURL, nil
	}

}

func (storage *InMemoryStorage) SelectAllUserRecords(ctx context.Context, userID string) ([]commontypes.UserRecordData, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return nil, nil
	}
}

func (storage *InMemoryStorage) ReadUserID(ctx context.Context, shortURLKey string) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		return "", nil
	}
}

func (storage *InMemoryStorage) MarkAsDeleted(ctx context.Context, shortURLKeys []string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
