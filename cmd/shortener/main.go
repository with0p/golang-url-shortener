package main

import (
	"fmt"
	"net/http"

	"github.com/with0p/golang-url-shortener.git/internal/app/initializer"
)

func main() {
	handler, config := initializer.InitWithInMemoryStorage()

	fmt.Println("Run on " + config.BaseURL)

	err := http.ListenAndServe(config.BaseURL, handler.GetHTTPHandler())
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
