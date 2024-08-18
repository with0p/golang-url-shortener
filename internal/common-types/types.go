package commontypes

type BatchRecord struct {
	ID          string
	ShortURLKey string
	ShortURL    string
	FullURL     string
}

type RecordToBatch struct {
	ID      string
	FullURL string
}
