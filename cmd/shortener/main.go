package main

import (
	"net/http"

	"github.com/with0p/golang-url-shortener.git/internal/app/initializer"
	"github.com/with0p/golang-url-shortener.git/internal/logger"
)

func main() {
	config := initializer.InitConfig()
	handler, initError := initializer.InitWithLocalFileStorage(config)

	if initError != nil {
		logger.LogError(initError)
		return
	}

	logger.LogInfo("Run on " + config.BaseURL)

	err := http.ListenAndServe(config.BaseURL, handler.GetHTTPHandler())
	if err != nil {
		logger.LogError(err)
		return
	}
}
