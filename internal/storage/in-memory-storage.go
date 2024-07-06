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

func (storage *InMemoryStorage) Write(key string, value string) error {
	storage.urlMap[key] = value
	return nil
}

func (storage *InMemoryStorage) Read(key string) (string, error) {
	value, ok := storage.urlMap[key]

	if !ok {
		return "", errors.New("Not found")
	}

	return value, nil
}

func (storage *InMemoryStorage) GetStorageSize() int {
	return len(storage.urlMap)
}
