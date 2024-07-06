package service

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"net/url"

	"github.com/with0p/golang-url-shortener.git/internal/storage"
)

type ShortURLService struct {
	storage      storage.Storage
	shortURLHost string
}

func NewShortURLService(currentStorage storage.Storage, shortURLHost string) *ShortURLService {
	return &ShortURLService{
		storage:      currentStorage,
		shortURLHost: shortURLHost,
	}
}

func (s *ShortURLService) MakeShortURL(trueURL string) (string, error) {
	_, urlParseError := url.ParseRequestURI(trueURL)

	if urlParseError != nil {
		return "", errors.New("not a URL")
	}

	shortURLId := GenerateShortURLId([]byte(trueURL))

	if err := s.storage.Write(shortURLId, trueURL); err != nil {
		return "", errors.New("could not make URL record")
	}

	return s.shortURLHost + "/" + shortURLId, nil
}

func (s *ShortURLService) GetTrueURL(id string) (string, error) {
	return s.storage.Read(id)
}

func GenerateShortURLId(fullURLByte []byte) string {
	hash := md5.New()
	hash.Write(fullURLByte)
	hashBytes := hash.Sum(nil)
	return hex.EncodeToString(hashBytes)[:8]
}
