package service

import (
	"context"

	commontypes "github.com/with0p/golang-url-shortener.git/internal/common-types"
)

type Service interface {
	MakeShortURL(ctx context.Context, trueURL string) (string, error)
	GetTrueURL(ctx context.Context, id string) (string, error)
	MakeShortURLBatch(ctx context.Context, recordsIn []commontypes.RecordToBatch) ([]commontypes.BatchRecord, error)
	GetAllUserRecords(ctx context.Context, userID string) ([]commontypes.UserRecord, error)
	DeleteUserURLs(userID string, shortURLKeys []string)
}
