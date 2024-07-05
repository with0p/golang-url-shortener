package main

import (
	"net/http"

	"github.com/with0p/golang-url-shortener.git/cmd/shortener/config"
	"github.com/with0p/golang-url-shortener.git/cmd/shortener/handlers"
)

func main() {
	config.ParseConfig()

	router := handlers.ServerRouter()

	err1 := http.ListenAndServe(config.Config.BaseURL, router)
	if err1 != nil {
		panic(err1)
	}

}
