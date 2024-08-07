package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/with0p/golang-url-shortener.git/internal/logger"
)

type DBStorage struct {
	db *sql.DB
}

func NewDBStorage(db *sql.DB) (*DBStorage, error) {
	err := initTable(db)
	if err != nil {
		return nil, err
	}
	return &DBStorage{db: db}, nil
}

func initTable(db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS shortener (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        full_url TEXT NOT NULL,
        short_url_key TEXT NOT NULL
    );`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := db.ExecContext(ctx, query)
	if err != nil {
		logger.LogError(err)
	}

	logger.LogInfo("Table %s created successfully")
	return err
}

func (storage *DBStorage) Write(shortURLKey string, fullURL string) error {
	queryExists := `
    SELECT EXISTS (
        SELECT 1 
        FROM shortener 
        WHERE short_url_key = $1
		);`

	ctxExists, cancelExists := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelExists()

	var exists bool
	errExists := storage.db.QueryRowContext(ctxExists, queryExists, shortURLKey).Scan(&exists)

	if errExists != nil {
		return errExists
	}

	if exists {
		return nil
	}

	ctxInsert, cancelInsert := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelInsert()

	queryInsert := `
    INSERT INTO shortener (full_url, short_url_key) 
    VALUES ($1, $2);`

	_, errInsert := storage.db.ExecContext(ctxInsert, queryInsert, fullURL, shortURLKey)
	return errInsert
}

func (storage *DBStorage) Read(shortURLKey string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
    SELECT full_url 
    FROM shortener 
    WHERE short_url_key = $1;`

	var fullURL string
	err := storage.db.QueryRowContext(ctx, query, shortURLKey).Scan(&fullURL)
	if err != nil {
		return "", err
	}
	return fullURL, nil
}
