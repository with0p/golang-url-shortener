package main

import (
	"net/http"

	"github.com/with0p/golang-url-shortener.git/cmd/shortener/handlers"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, handlers.URLShortener)
	mux.HandleFunc(`/{id}`, handlers.GetTrueURL)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}

}
