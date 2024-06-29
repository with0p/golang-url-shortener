package main

import (
	"net/http"

	handlers "github.com/with0p/golang-url-shortener.git/internal/app/handlers"
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
