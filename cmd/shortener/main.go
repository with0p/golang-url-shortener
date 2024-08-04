package main

import (
	"net/http"

	"github.com/with0p/golang-url-shortener.git/internal/app/initializer"
	"github.com/with0p/golang-url-shortener.git/internal/logger"

	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	config := initializer.InitConfig()
	handler, initError := initializer.InitWithLocalFileStorage(config)

	if initError != nil {
		logger.LogError(initError)
		return
	}

	db, dbErr := sql.Open("pgx", config.DataBaseAddress)
	if dbErr != nil {
		logger.LogError(dbErr)
	}
	defer db.Close()

	logger.LogInfo("Run on " + config.BaseURL)

	err := http.ListenAndServe(config.BaseURL, handler.GetHTTPHandler(db))
	if err != nil {
		logger.LogError(err)
		return
	}
}
