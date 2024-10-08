package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	commontypes "github.com/with0p/golang-url-shortener.git/internal/common-types"
	customerrors "github.com/with0p/golang-url-shortener.git/internal/custom-errors"
)

type DBStorage struct {
	db *sql.DB
}

func NewDBStorage(ctx context.Context, db *sql.DB) (*DBStorage, error) {
	err := initTable(ctx, db)
	if err != nil {
		return nil, err
	}
	return &DBStorage{db: db}, nil
}

func initTable(ctx context.Context, db *sql.DB) error {
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

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return tr.Commit()
	}
}

func (storage *DBStorage) Read(ctx context.Context, shortURLKey string) (string, error) {
	query := `
	SELECT full_url 
	FROM shortener 
	WHERE short_url_key = $1;`

	var fullURL string
	err := storage.db.QueryRowContext(ctx, query, shortURLKey).Scan(&fullURL)
	if err != nil {
		return "", err
	}

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		return fullURL, nil
	}
}

func (storage *DBStorage) Write(ctx context.Context, shortURLKey string, fullURL string) error {
	queryInsert := `
    INSERT INTO shortener (full_url, short_url_key) 
    VALUES ($1, $2);`

	_, errInsert := storage.db.ExecContext(ctx, queryInsert, fullURL, shortURLKey)
	if errInsert != nil {
		var pgErr *pgconn.PgError
		if errors.As(errInsert, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			errInsert = customerrors.ErrUniqueKeyConstrantViolation
		}
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return errInsert
	}
}

func (storage *DBStorage) WriteBatch(ctx context.Context, records []commontypes.BatchRecord) error {
	tr, err := storage.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	for _, r := range records {
		queryInsert := `
    INSERT INTO shortener (full_url, short_url_key) 
    VALUES ($1, $2);`

		_, errInsert := tr.ExecContext(ctx, queryInsert, r.FullURL, r.ShortURLKey)
		if errInsert != nil {
			tr.Rollback()
			return errInsert
		}
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return tr.Commit()
	}
}
