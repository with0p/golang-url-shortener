package config

var MockConfiguration = &Config{
	BaseURL:         "localhost:8080",
	ShortURL:        "http://localhost:8080",
	FileStoragePath: "internal/storage/local-file/local-storage.json",
	DataBaseAddress: "host=localhost port=5435 user=postgres password=1234 dbname=postgres sslmode=disable",
}
