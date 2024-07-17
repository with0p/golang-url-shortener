package initializer

import (
	"github.com/with0p/golang-url-shortener.git/internal/config"
	"github.com/with0p/golang-url-shortener.git/internal/handler"
	"github.com/with0p/golang-url-shortener.git/internal/service"
	"github.com/with0p/golang-url-shortener.git/internal/storage"
)

func runInit(storage storage.Storage) (*handler.URLHandler, *config.Config) {
	config := config.GetConfig()
	service := service.NewShortURLService(storage, config.ShortURL)
	urlHandler := handler.NewURLHandler(service)

	return urlHandler, config
}
