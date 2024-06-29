package storage

type URLStorage map[string]string

var URLMapStorage URLStorage
var URLMap *URLStorage

func (storage *URLStorage) Set(key string, value string) {
	(*storage)[key] = value
}

func (storage *URLStorage) Get(key string) (string, bool) {
	value, ok := (*storage)[key]
	return value, ok
}

func GetURLMap() *URLStorage {
	if URLMapStorage == nil {
		InitMap()
	}
	return &URLMapStorage
}

func InitMap() {
	URLMapStorage = make(map[string]string)
}

// type Storage struct {
// 	data map[string]string
// }

// func NewStorage(initialData map[string]string) *Storage {
// 	return &Storage{
// 		data: initialData,
// 	}
// }

// func (s *Storage) GetItem(key string) string {
// 	return s.data[key]
// }

// func (s *Storage) SetItem(key, value string) {
// 	s.data[key] = value
// }
