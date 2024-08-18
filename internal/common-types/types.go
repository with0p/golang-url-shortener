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

type UserRecordData struct {
	ShortURLKey string
	FullURL     string
}

type UserRecord struct {
	ShortURL string `json:"short_url"`
	FullURL  string `json:"original_url"`
}
