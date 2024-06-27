package main

import (
	"io"
	"net/http"
	"strconv"
	"strings"
)

var UrlMap map[string]string

func getUrlMapKey() string {
	return "shorturl" + strconv.Itoa(len(UrlMap))
}

func UrlShortener(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Not a POST requests", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(req.Body)

	if err != nil {
		res.Write([]byte(err.Error()))
		return
	}

	urlKey := getUrlMapKey()
	UrlMap[urlKey] = string(body)

	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(201)
	res.Write([]byte("http://localhost:8080/" + urlKey))
}

func GetTrueURL(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Not a GET requests", http.StatusMethodNotAllowed)
		return
	}

	pathSplitted := strings.Split(req.URL.Path, "/")
	id := pathSplitted[len(pathSplitted)-1]

	trueURL, ok := UrlMap[id]

	if !ok {
		http.Error(res, "Not found", http.StatusNotFound)
		return
	}

	http.Redirect(res, req, trueURL, http.StatusBadRequest)
}

func main() {
	UrlMap = make(map[string]string)

	mux := http.NewServeMux()
	mux.HandleFunc(`/`, UrlShortener)
	mux.HandleFunc(`/{id}`, GetTrueURL)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}

}
