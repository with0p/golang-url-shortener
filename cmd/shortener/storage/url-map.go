package storage

type URLStorageMap map[string]string
type URLStorage struct {
	urlMap URLStorageMap
}

var URLMap *URLStorage

func (storage *URLStorage) Set(key string, value string) {
	storage.urlMap[key] = value
}

func (storage *URLStorage) Get(key string) (string, bool) {
	value, ok := storage.urlMap[key]
	return value, ok
}

func (storage *URLStorage) GetStorageSize() int {
	return len(storage.urlMap)
}

func GetURLMap() *URLStorage {
	if URLMap == nil {
		InitMap()
	}
	return URLMap
}

func InitMap() {
	URLMap = &URLStorage{
		urlMap: map[string]string{},
	}
}
