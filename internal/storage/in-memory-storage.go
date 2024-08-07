package storage

import "errors"

type URLStorageMap map[string]string

type InMemoryStorage struct {
	urlMap URLStorageMap
}

func NewInMemoryStorage(storageMap URLStorageMap) *InMemoryStorage {
	return &InMemoryStorage{
		urlMap: storageMap,
	}
}

func (storage *InMemoryStorage) Write(shortURLKey string, fullURL string) error {
	storage.urlMap[shortURLKey] = fullURL
	return nil
}

func (storage *InMemoryStorage) Read(shortURLKey string) (string, error) {
	fullURL, ok := storage.urlMap[shortURLKey]

	if !ok {
		return "", errors.New("not found")
	}

	return fullURL, nil
}

func (storage *InMemoryStorage) GetStorageSize() int {
	return len(storage.urlMap)
}
