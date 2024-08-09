package service

import commontypes "github.com/with0p/golang-url-shortener.git/internal/common-types"

type Service interface {
	MakeShortURL(trueURL string) (string, error)
	GetTrueURL(id string) (string, error)
	MakeShortURLBatch(recordsIn *[]commontypes.RecordToBatch) (*[]commontypes.BatchRecord, error)
}
