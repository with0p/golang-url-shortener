package config

import (
	"flag"
	"net/url"
	"os"

	"github.com/with0p/golang-url-shortener.git/internal/logger"
)

// const defaultFileStoragePath = "internal/storage/local-file/local-storage.json"
// const defaultDataBaseAddress = "host=localhost port=5435 user=postgres password=1234 dbname=postgres sslmode=disable"
const defaultHost = "localhost"
const defaultPort = "8080"
const defaultFileStoragePath = ""
const defaultDataBaseAddress = ""

type Config struct {
	BaseURL         string
	ShortURL        string
	FileStoragePath string
	DataBaseAddress string
}

var configuration *Config

func GetConfig() *Config {
	if configuration == nil {
		var conf = Config{}

		flag.StringVar(&conf.BaseURL, "a", defaultHost+":"+defaultPort, "base URL")
		flag.StringVar(&conf.ShortURL, "b", "http://"+defaultHost+":"+defaultPort, "short URL")
		flag.StringVar(&conf.FileStoragePath, "f", defaultFileStoragePath, "storage path")
		flag.StringVar(&conf.DataBaseAddress, "d", defaultDataBaseAddress, "database address")
		flag.Parse()

		if envServerAddress := os.Getenv("SERVER_ADDRESS"); envServerAddress != "" {
			conf.BaseURL = envServerAddress
		}

		if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
			conf.ShortURL = envBaseURL
		}

		if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
			conf.FileStoragePath = envFileStoragePath
		}

		if envDatabaseStoragePath := os.Getenv("DATABASE_DSN"); envDatabaseStoragePath != "" {
			conf.DataBaseAddress = envDatabaseStoragePath
		}

		configuration = &Config{
			BaseURL:         URLParseHelper(conf.BaseURL),
			ShortURL:        "http://" + URLParseHelper(conf.ShortURL),
			FileStoragePath: conf.FileStoragePath,
			DataBaseAddress: conf.DataBaseAddress,
		}
	}

	return configuration

}

func URLParseHelper(str string) string {
	parsedURL, err := url.ParseRequestURI(str)
	if err != nil {
		logger.LogError(err)
	}

	if parsedURL.Host == "" {
		parsedURL, err = url.ParseRequestURI("http://" + str)
		if err != nil {
			logger.LogError(err)
		}
	}

	port := parsedURL.Port()
	if port == "" {
		port = defaultPort
	}

	host := parsedURL.Hostname()
	if host == "" {
		host = defaultHost
	}

	return host + ":" + port
}
