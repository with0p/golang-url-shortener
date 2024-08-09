package storage

import commontypes "github.com/with0p/golang-url-shortener.git/internal/common-types"

type Storage interface {
	Read(shortURLKey string) (string, error)
	Write(shortURLKey string, fullURL string) error
	WriteBatch(records *[]commontypes.BatchRecord) error
}
