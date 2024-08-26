package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"net/url"
	"sync"
	"time"

	"github.com/with0p/golang-url-shortener.git/internal/auth"
	commontypes "github.com/with0p/golang-url-shortener.git/internal/common-types"
	customerrors "github.com/with0p/golang-url-shortener.git/internal/custom-errors"
	"github.com/with0p/golang-url-shortener.git/internal/logger"
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
	userID, err := auth.GetUserIDFromCtx(ctx)

	if err != nil {
		return "", err
	}

	if err := s.storage.Write(ctx, userID, shortURLId, trueURL); err != nil {
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

	userID, err := auth.GetUserIDFromCtx(ctx)

	if err != nil {
		return nil, err
	}

	if err := s.storage.WriteBatch(ctx, userID, batchData); err != nil {
		return nil, errors.New("could not make Batch URL record")
	}

	return batchData, nil
}

func (s *ShortURLService) GetAllUserRecords(ctx context.Context, userID string) ([]commontypes.UserRecord, error) {
	recordData, err := s.storage.SelectAllUserRecords(ctx, userID)

	if err != nil {
		return nil, err
	}

	userRecords := make([]commontypes.UserRecord, len(recordData))

	for i, data := range recordData {
		userRecords[i] = commontypes.UserRecord{
			ShortURL: s.shortURLHost + "/" + data.ShortURLKey,
			FullURL:  data.FullURL,
		}
	}

	return userRecords, nil
}

func (s ShortURLService) DeleteUserURLs(userID string, shortURLKeys []string) {
	ctx1, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	inCh := make(chan string, len(shortURLKeys))
	resultCh := make(chan string)

	go func() {
		for _, id := range shortURLKeys {
			inCh <- id
		}
		close(inCh)
	}()

	var wg sync.WaitGroup

	for w := 1; w <= 3; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for shortURLKey := range inCh {
				recordUserID, err := s.storage.ReadUserID(ctx1, shortURLKey)

				if err == nil && recordUserID == userID {
					resultCh <- shortURLKey
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	var results []string
	for res := range resultCh {
		results = append(results, res)
	}

	if err := s.storage.MarkAsDeleted(ctx1, results); err != nil {
		logger.LogError(err)
	}

}

func generateShortURLId(fullURLByte []byte) string {
	hash := md5.New()
	hash.Write(fullURLByte)
	hashBytes := hash.Sum(nil)
	return hex.EncodeToString(hashBytes)[:8]
}
