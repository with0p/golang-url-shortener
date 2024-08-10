package service

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"net/url"

	"github.com/with0p/golang-url-shortener.git/internal/common-types"
	customerrors "github.com/with0p/golang-url-shortener.git/internal/custom-errors"
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

func (s *ShortURLService) GetTrueURL(id string) (string, error) {
	return s.storage.Read(id)
}

func (s *ShortURLService) MakeShortURL(trueURL string) (string, error) {
	_, urlParseError := url.ParseRequestURI(trueURL)

	if urlParseError != nil {
		return "", errors.New("not a URL")
	}

	shortURLId := generateShortURLId([]byte(trueURL))

	if err := s.storage.Write(shortURLId, trueURL); err != nil {
		if errors.Is(err, customerrors.ErrUniqueKeyConstrantViolation) {
			return s.shortURLHost + "/" + shortURLId, err
		}
		return "", err
	}

	return s.shortURLHost + "/" + shortURLId, nil
}

func (s *ShortURLService) MakeShortURLBatch(recordsIn *[]commontypes.RecordToBatch) (*[]commontypes.BatchRecord, error) {
	var batchData []commontypes.BatchRecord

	for _, reqRec := range *recordsIn {
		_, urlParseError := url.ParseRequestURI(reqRec.FullURL)

		if urlParseError != nil {
			continue
		}

		shortURLId := generateShortURLId([]byte(reqRec.FullURL))

		batchData = append(batchData, commontypes.BatchRecord{
			ID:          reqRec.ID,
			ShortURLKey: shortURLId,
			ShortURL:    s.shortURLHost + "/" + shortURLId,
			FullURL:     reqRec.FullURL,
		})
	}

	if err := s.storage.WriteBatch(&batchData); err != nil {
		return nil, errors.New("could not make Batch URL record")
	}

	return &batchData, nil
}

func generateShortURLId(fullURLByte []byte) string {
	hash := md5.New()
	hash.Write(fullURLByte)
	hashBytes := hash.Sum(nil)
	return hex.EncodeToString(hashBytes)[:8]
}
