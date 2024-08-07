package storage

type Storage interface {
	Write(shortURLKey string, fullURL string) error
	Read(shortURLKey string) (string, error)
}
