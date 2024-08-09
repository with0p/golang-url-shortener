package main

import (
	"net/http"

	"github.com/with0p/golang-url-shortener.git/internal/app/initializer"
	"github.com/with0p/golang-url-shortener.git/internal/handler"
	"github.com/with0p/golang-url-shortener.git/internal/logger"

	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	config := initializer.InitConfig()
	var dataBase *sql.DB
	var handler *handler.URLHandler
	var initError error

	if config.DataBaseAddress != "" {
		db, dbErr := sql.Open("pgx", config.DataBaseAddress)
		if dbErr != nil {
			logger.LogError(dbErr)
		}
		defer db.Close()
		dataBase = db

		handler, initError = initializer.InitWithDBStorage(config, dataBase)
	} else if config.FileStoragePath != "" {
		handler, initError = initializer.InitWithLocalFileStorage(config)
	} else {
		handler, initError = initializer.InitWithInMemoryStorage(config)
	}

	if initError != nil {
		logger.LogError(initError)
		return
	}

	logger.LogInfo("Run on " + config.BaseURL)

	err := http.ListenAndServe(config.BaseURL, handler.GetHTTPHandler(dataBase))
	if err != nil {
		logger.LogError(err)
		return
	}
}
