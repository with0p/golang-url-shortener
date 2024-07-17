package service

type Service interface {
	MakeShortURL(trueURL string) (string, error)
	GetTrueURL(id string) (string, error)
}
