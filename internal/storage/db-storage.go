package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	commontypes "github.com/with0p/golang-url-shortener.git/internal/common-types"
	customerrors "github.com/with0p/golang-url-shortener.git/internal/custom-errors"
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
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tr, errTr := db.BeginTx(ctx, nil)
	if errTr != nil {
		return errTr
	}

	query := `
    CREATE TABLE IF NOT EXISTS shortener (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        full_url TEXT NOT NULL,
        short_url_key TEXT NOT NULL
    );`

	tr.ExecContext(ctx, query)

	tr.ExecContext(ctx, `CREATE UNIQUE INDEX IF NOT EXISTS short_url_key_index ON shortener (short_url_key)`)

	return tr.Commit()
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

func (storage *DBStorage) Write(shortURLKey string, fullURL string) error {
	ctxInsert, cancelInsert := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelInsert()

	queryInsert := `
    INSERT INTO shortener (full_url, short_url_key) 
    VALUES ($1, $2);`

	_, errInsert := storage.db.ExecContext(ctxInsert, queryInsert, fullURL, shortURLKey)
	if errInsert != nil {
		var pgErr *pgconn.PgError
		if errors.As(errInsert, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			errInsert = customerrors.ErrUniqueKeyConstrantViolation
		}
	}

	return errInsert
}

func (storage *DBStorage) WriteBatch(records *[]commontypes.BatchRecord) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tr, err := storage.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	for _, r := range *records {
		ctxInsert, cancelInsert := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancelInsert()

		queryInsert := `
    INSERT INTO shortener (full_url, short_url_key) 
    VALUES ($1, $2);`

		_, errInsert := tr.ExecContext(ctxInsert, queryInsert, r.FullURL, r.ShortURLKey)
		if errInsert != nil {
			tr.Rollback()
			return errInsert
		}
	}

	return tr.Commit()
}
