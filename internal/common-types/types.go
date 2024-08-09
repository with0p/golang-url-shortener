package commontypes

type BatchRecord struct {
	ShortURLKey string
	FullURL     string
}

type RecordToBatch struct {
	Id          string
	OriginalURL string
}
