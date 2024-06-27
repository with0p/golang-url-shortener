package main

import (
	"io"
	"net/http"
	"strconv"
	"strings"
)

var URLMap map[string]string

func getURLMapKey() string {
	return "shorturl" + strconv.Itoa(len(URLMap))
}

func URLShortener(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Not a POST requests", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(req.Body)

	if err != nil {
		res.Write([]byte(err.Error()))
		return
	}

	urlKey := getURLMapKey()
	URLMap[urlKey] = string(body)

	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(201)
	res.Write([]byte("http://localhost:8080/" + urlKey))
}

func GetTrueURL(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Not a GET requests", http.StatusBadRequest)
		return
	}

	pathSplitted := strings.Split(req.URL.Path, "/")
	id := pathSplitted[len(pathSplitted)-1]

	trueURL, ok := URLMap[id]

	if !ok {
		http.Error(res, "Not found", http.StatusNotFound)
		return
	}

	http.Redirect(res, req, trueURL, http.StatusTemporaryRedirect)
}

func main() {
	URLMap = make(map[string]string)

	mux := http.NewServeMux()
	mux.HandleFunc(`/`, URLShortener)
	mux.HandleFunc(`/{id}`, GetTrueURL)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}

}
