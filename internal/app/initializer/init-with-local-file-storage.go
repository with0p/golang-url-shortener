package initializer

import (
	"errors"
	"fmt"

	"github.com/with0p/golang-url-shortener.git/internal/config"
	"github.com/with0p/golang-url-shortener.git/internal/handler"
	"github.com/with0p/golang-url-shortener.git/internal/logger"
	"github.com/with0p/golang-url-shortener.git/internal/storage/local-file"
)

func InitWithLocalFileStorage(config *config.Config) (*handler.URLHandler, error) {

	storage, err := storage.NewLocalFileStorage(config.FileStoragePath)
	if err != nil {
		logger.LogError(err)
		return nil, errors.New("cannot init local file storage")
	}
	logger.LogInfo(fmt.Sprintf(`File storage path: %s`, config.FileStoragePath))
	return runInit(storage, config), nil
}
