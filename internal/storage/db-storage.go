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
		user_id TEXT NOT NULL,
        full_url TEXT NOT NULL,
        short_url_key TEXT NOT NULL,
		is_deleted BOOL DEFAULT FALSE
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
	SELECT full_url, is_deleted 
	FROM shortener 
	WHERE short_url_key = $1;`

	var fullURL string
	var isDeleted bool
	err := storage.db.QueryRowContext(ctx, query, shortURLKey).Scan(&fullURL, &isDeleted)
	if err != nil {
		return "", err
	} else if isDeleted {
		return "", customerrors.ErrRecordHasBeenDeleted
	}

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		return fullURL, nil
	}
}

func (storage *DBStorage) ReadUserID(ctx context.Context, shortURLKey string) (string, error) {
	query := `
	SELECT user_id 
	FROM shortener 
	WHERE short_url_key = $1;`

	var userID string
	err := storage.db.QueryRowContext(ctx, query, shortURLKey).Scan(&userID)
	if err != nil {
		return "", err
	}
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		return userID, nil
	}
}

func (storage *DBStorage) Write(ctx context.Context, userID string, shortURLKey string, fullURL string) error {
	queryInsert := `
    INSERT INTO shortener (user_id, full_url, short_url_key) 
    VALUES ($1, $2, $3);`

	_, errInsert := storage.db.ExecContext(ctx, queryInsert, userID, fullURL, shortURLKey)
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

func (storage *DBStorage) WriteBatch(ctx context.Context, userID string, records []commontypes.BatchRecord) error {
	tr, err := storage.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	for _, r := range records {
		queryInsert := `
    INSERT INTO shortener (user_id, full_url, short_url_key) 
    VALUES ($1, $2, $3);`

		_, errInsert := tr.ExecContext(ctx, queryInsert, userID, r.FullURL, r.ShortURLKey)
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

func (storage *DBStorage) SelectAllUserRecords(ctx context.Context, userID string) ([]commontypes.UserRecordData, error) {
	query := `
	SELECT full_url, short_url_key 
	FROM shortener 
	WHERE user_id = $1;`

	rows, err := storage.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userRecordsData []commontypes.UserRecordData

	for rows.Next() {
		var record commontypes.UserRecordData
		err := rows.Scan(&record.FullURL, &record.ShortURLKey)
		if err != nil {
			continue
		}
		userRecordsData = append(userRecordsData, record)

	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return userRecordsData, nil
	}
}

func (storage *DBStorage) MarkAsDeleted(ctx context.Context, shortURLKeys []string) error {
	tr, err := storage.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	for _, key := range shortURLKeys {

		queryUpdate := `
    	UPDATE shortener
		SET is_deleted = true
		WHERE short_url_key = $1;`

		_, errUpdate := tr.ExecContext(ctx, queryUpdate, key)
		if errUpdate != nil {
			tr.Rollback()
			return errUpdate
		}
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return tr.Commit()
	}
}
