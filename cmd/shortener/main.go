package main

import (
	"net/http"

	"github.com/with0p/golang-url-shortener.git/cmd/shortener/handlers"
)

func main() {
	err := http.ListenAndServe(":8080", handlers.ServerRouter())
	if err != nil {
		panic(err)
	}
}
