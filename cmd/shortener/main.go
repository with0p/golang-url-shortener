package main

import (
	"net/http"

	"github.com/with0p/golang-url-shortener.git/cmd/shortener/config"
	"github.com/with0p/golang-url-shortener.git/cmd/shortener/handlers"
)

func main() {
	config.ParseCMDFlags()

	router := handlers.ServerRouter()

	// go func() {
	// 	err2 := http.ListenAndServe(config.CMDFlags.ShortURL, router)
	// 	if err2 != nil {
	// 		panic(err2)
	// 	}
	// }()

	err1 := http.ListenAndServe(config.CMDFlags.BaseURL, router)
	if err1 != nil {
		panic(err1)
	}

}
