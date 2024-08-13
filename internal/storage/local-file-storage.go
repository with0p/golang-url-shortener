package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"os"

	commontypes "github.com/with0p/golang-url-shortener.git/internal/common-types"
	"github.com/with0p/golang-url-shortener.git/internal/logger"
	localfile "github.com/with0p/golang-url-shortener.git/internal/storage/local-file"
)

type LocalFileStorage struct {
	filePath string
}

func NewLocalFileStorage(filePath string) (*LocalFileStorage, error) {
	return &LocalFileStorage{filePath: filePath}, nil
}

func (storage *LocalFileStorage) Write(ctx context.Context, shortURLKey string, fullURL string) error {
	file, err := os.OpenFile(storage.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logger.LogError(err)
		return err
	}
	defer file.Close()

	data, err := json.Marshal(localfile.NewLocalFileRecord(shortURLKey, fullURL))
	if err != nil {
		logger.LogError(err)
		return err
	}
	data = append(data, '\n')

	_, err = file.Write(data)

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return err
	}
}

func (storage *LocalFileStorage) WriteBatch(ctx context.Context, records []commontypes.BatchRecord) error {
	file, err := os.OpenFile(storage.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logger.LogError(err)
		return err
	}
	defer file.Close()

	var dataToWrite []byte

	for _, r := range records {
		data, err := json.Marshal(localfile.NewLocalFileRecord(r.ShortURLKey, r.FullURL))
		if err != nil {
			logger.LogError(err)
			return err
		}
		data = append(data, '\n')
		dataToWrite = append(dataToWrite, data...)
	}

	_, err = file.Write(dataToWrite)

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return err
	}
}

func (storage *LocalFileStorage) Read(ctx context.Context, shortURLKey string) (string, error) {
	file, err := os.OpenFile(storage.filePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.LogError(err)
		return "", err
	}

	fileData, err := readFileToMap(file)
	if err != nil {
		logger.LogError(err)
		return "", err
	}
	defer file.Close()

	fullURL, ok := fileData[shortURLKey]

	if !ok {
		return "", errors.New("not found")
	}

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		return fullURL, nil
	}
}

func readFileToMap(file *os.File) (map[string]string, error) {
	fileData := map[string]string{}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		data := scanner.Bytes()

		record := localfile.LocalFileRecord{}

		err := json.Unmarshal(data, &record)
		if err != nil {
			logger.LogError(err)
			return nil, err
		}

		fileData[record.ShortURL] = record.OriginalURL
	}

	return fileData, nil
}
