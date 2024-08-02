package initializer

import (
	"github.com/with0p/golang-url-shortener.git/internal/config"
	"github.com/with0p/golang-url-shortener.git/internal/handler"
	"github.com/with0p/golang-url-shortener.git/internal/service"
	"github.com/with0p/golang-url-shortener.git/internal/storage"
)

func InitConfig() *config.Config {
	return config.GetConfig()
}

func runInit(storage storage.Storage, config *config.Config) *handler.URLHandler {
	service := service.NewShortURLService(storage, config.ShortURL)
	urlHandler := handler.NewURLHandler(service)

	return urlHandler
}
