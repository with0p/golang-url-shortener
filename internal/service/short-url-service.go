package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"net/url"

	commontypes "github.com/with0p/golang-url-shortener.git/internal/common-types"
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

func (s *ShortURLService) GetTrueURL(ctx context.Context, id string) (string, error) {
	return s.storage.Read(ctx, id)
}

func (s *ShortURLService) MakeShortURL(ctx context.Context, trueURL string) (string, error) {
	_, urlParseError := url.ParseRequestURI(trueURL)

	if urlParseError != nil {
		return "", errors.New("not a URL")
	}

	shortURLId := generateShortURLId([]byte(trueURL))

	if err := s.storage.Write(ctx, shortURLId, trueURL); err != nil {
		if errors.Is(err, customerrors.ErrUniqueKeyConstrantViolation) {
			return s.shortURLHost + "/" + shortURLId, err
		}
		return "", errors.New("could not make URL record")
	}

	return s.shortURLHost + "/" + shortURLId, nil
}

func (s *ShortURLService) MakeShortURLBatch(ctx context.Context, recordsIn []commontypes.RecordToBatch) ([]commontypes.BatchRecord, error) {
	batchData := make([]commontypes.BatchRecord, len(recordsIn))

	for i, reqRec := range recordsIn {
		_, urlParseError := url.ParseRequestURI(reqRec.FullURL)

		if urlParseError != nil {
			continue
		}

		shortURLId := generateShortURLId([]byte(reqRec.FullURL))

		batchData[i] = commontypes.BatchRecord{
			ID:          reqRec.ID,
			ShortURLKey: shortURLId,
			ShortURL:    s.shortURLHost + "/" + shortURLId,
			FullURL:     reqRec.FullURL,
		}
	}

	if err := s.storage.WriteBatch(ctx, batchData); err != nil {
		return nil, errors.New("could not make Batch URL record")
	}

	return batchData, nil
}

func generateShortURLId(fullURLByte []byte) string {
	hash := md5.New()
	hash.Write(fullURLByte)
	hashBytes := hash.Sum(nil)
	return hex.EncodeToString(hashBytes)[:8]
}
