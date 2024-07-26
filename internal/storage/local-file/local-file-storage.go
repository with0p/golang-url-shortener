package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"

	"github.com/with0p/golang-url-shortener.git/internal/logger"
)

type LocalFileStorage struct {
	filePath string
}

func NewLocalFileStorage(filePath string) (*LocalFileStorage, error) {
	return &LocalFileStorage{filePath: filePath}, nil
}

func (storage *LocalFileStorage) Write(key string, value string) error {
	file, err := os.OpenFile(storage.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logger.LogError(err)
		return err
	}
	defer file.Close()

	data, err := json.Marshal(NewLocalFileRecord(key, value))
	if err != nil {
		logger.LogError(err)
		return err
	}
	data = append(data, '\n')

	_, err = file.Write(data)

	return err
}

func (storage *LocalFileStorage) Read(key string) (string, error) {
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

	value, ok := fileData[key]

	if !ok {
		return "", errors.New("not found")
	}

	return value, nil
}

func readFileToMap(file *os.File) (map[string]string, error) {
	fileData := map[string]string{}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		data := scanner.Bytes()

		record := LocalFileRecord{}

		err := json.Unmarshal(data, &record)
		if err != nil {
			logger.LogError(err)
			return nil, err
		}

		fileData[record.ShortURL] = record.OriginalURL
	}

	return fileData, nil
}
