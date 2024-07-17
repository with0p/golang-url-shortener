package config

import (
	"flag"
	"fmt"
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
		var conf = Config{}

		flag.StringVar(&conf.BaseURL, "a", defaultHost+":"+defaultPort, "base URL")
		flag.StringVar(&conf.ShortURL, "b", "http://"+defaultHost+":"+defaultPort, "short URL")
		flag.Parse()

		if envServerAddress := os.Getenv("SERVER_ADDRESS"); envServerAddress != "" {
			conf.BaseURL = envServerAddress
		}

		if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
			conf.ShortURL = envBaseURL
		}

		configuration = &Config{
			BaseURL:  URLParseHelper(conf.BaseURL),
			ShortURL: "http://" + URLParseHelper(conf.ShortURL),
		}
	}

	return configuration

}

func URLParseHelper(str string) string {
	parsedURL, err := url.ParseRequestURI(str)
	if err != nil {
		fmt.Println(err.Error())
	}

	if parsedURL.Host == "" {
		parsedURL, err = url.ParseRequestURI("http://" + str)
		if err != nil {
			fmt.Println(err.Error())
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
