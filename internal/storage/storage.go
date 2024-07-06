package storage

type Storage interface {
	Write(key string, value string) error
	Read(key string) (string, error)
}
