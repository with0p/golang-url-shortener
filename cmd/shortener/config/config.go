package config

import (
	"flag"
	"net/url"
)

var CMDFlags struct {
	BaseURL  string
	ShortURL string
}

const defaultHOST = "localhost"
const defaultBasePort = "8080"
const defaultShortPort = "8888"

func ParseCMDFlags() {
	var BaseURL string
	var ShortURL string

	flag.StringVar(&BaseURL, "a", "http://localhost:8080", "base URL")
	flag.StringVar(&ShortURL, "b", "http://localhost:8888", "short URL")
	flag.Parse()

	CMDFlags.BaseURL = URLParseHelper(BaseURL, defaultBasePort)
	CMDFlags.ShortURL = URLParseHelper(ShortURL, defaultShortPort)
}

func URLParseHelper(str string, defaultPort string) string {
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
		host = defaultHOST
	}

	return host + ":" + port
}
