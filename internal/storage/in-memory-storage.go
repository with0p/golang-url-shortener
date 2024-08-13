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

func (storage *InMemoryStorage) Write(ctx context.Context, shortURLKey string, fullURL string) error {
	storage.urlMap[shortURLKey] = fullURL
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func (storage *InMemoryStorage) WriteBatch(ctx context.Context, records []commontypes.BatchRecord) error {
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

func (storage *InMemoryStorage) GetStorageSize() int {
	return len(storage.urlMap)
}
