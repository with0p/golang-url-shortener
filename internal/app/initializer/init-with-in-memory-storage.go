package initializer

import (
	"github.com/with0p/golang-url-shortener.git/internal/config"
	"github.com/with0p/golang-url-shortener.git/internal/handler"
	"github.com/with0p/golang-url-shortener.git/internal/storage"
)

func InitWithInMemoryStorage(config *config.Config) (*handler.URLHandler, error) {
	inMemoryStorage := storage.NewInMemoryStorage(map[string]string{})
	return runInit(inMemoryStorage, config), nil
}
