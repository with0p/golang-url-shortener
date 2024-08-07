package initializer

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/with0p/golang-url-shortener.git/internal/config"
	"github.com/with0p/golang-url-shortener.git/internal/handler"
	"github.com/with0p/golang-url-shortener.git/internal/logger"
	"github.com/with0p/golang-url-shortener.git/internal/storage"
)

func InitWithDBStorage(config *config.Config, db *sql.DB) (*handler.URLHandler, error) {

	storage, err := storage.NewDBStorage(db)
	if err != nil {
		logger.LogError(err)
		return nil, errors.New("cannot init db storage")
	}
	logger.LogInfo(fmt.Sprintf(`DB address: %s`, config.DataBaseAddress))
	return runInit(storage, config), nil
}
