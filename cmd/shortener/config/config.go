package config

import (
	"flag"
)

var CMDFlags struct {
	BaseURL  string
	ShortURL string
}

func ParseCMDFlags() {
	flag.StringVar(&CMDFlags.BaseURL, "a", "localhost:8888", "base URL")
	flag.StringVar(&CMDFlags.ShortURL, "b", "localhost:8000", "short URL")
	flag.Parse()
}
