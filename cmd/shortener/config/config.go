package config

import (
	"flag"
	"net/url"
	"os"
)

var Config struct {
	BaseURL  string
	ShortURL string
}

const defaultHost = "localhost"
const defaultPort = "8080"

func ParseConfig() {
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

	Config.BaseURL = URLParseHelper(BaseURL)
	Config.ShortURL = "http://" + URLParseHelper(ShortURL)

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
