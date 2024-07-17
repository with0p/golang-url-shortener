package config

import (
	"flag"
	"net/url"
	"os"
)

const defaultHost = "localhost"
const defaultPort = "8080"

type Config struct {
	BaseURL  string
	ShortURL string
}

var configuration *Config

func GetConfig() *Config {
	if configuration == nil {
		var BaseURL string
		var ShortURL string

		flag.StringVar(&BaseURL, "a", defaultHost+":"+defaultPort, "base URL")
		flag.StringVar(&ShortURL, "b", "http://"+defaultHost+":"+defaultPort, "short URL")
		flag.Parse()

		if envServerAddress := os.Getenv("SERVER_ADDRESS"); envServerAddress != "" {
			BaseURL = envServerAddress
		}

		if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
			ShortURL = envBaseURL
		}

		configuration = &Config{
			BaseURL:  URLParseHelper(BaseURL),
			ShortURL: "http://" + URLParseHelper(ShortURL),
		}
	}

	return configuration

}

func URLParseHelper(str string) string {
	parsedURL, _ := url.ParseRequestURI(str)

	if parsedURL.Host == "" {
		parsedURL, _ = url.ParseRequestURI("http://" + str)
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
