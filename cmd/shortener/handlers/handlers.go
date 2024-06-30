package handlers

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/with0p/golang-url-shortener.git/cmd/shortener/storage"
)

func URLShortener(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Not a POST requests", http.StatusBadRequest)
		return
	}

	if req.Header.Get("content-type") != "text/plain" {
		http.Error(res, "Wrong content type", http.StatusBadRequest)
		return
	}

	defer req.Body.Close()
	body, bodyReadError := io.ReadAll(req.Body)
	if bodyReadError != nil {
		http.Error(res, bodyReadError.Error(), http.StatusBadRequest)
		return
	}

	_, urlParseError := url.ParseRequestURI((string(body)))
	if urlParseError != nil {
		http.Error(res, urlParseError.Error(), http.StatusBadRequest)
		return
	}

	urlKey := GenerateShortURL(body)

	storageInstance := storage.GetURLMap()
	storageInstance.Set(urlKey, string(body))

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

	trueURL, ok := storage.GetURLMap().Get(id)

	if !ok {
		http.Error(res, "Not found", http.StatusNotFound)
		return
	}

	http.Redirect(res, req, trueURL, http.StatusTemporaryRedirect)
}

func GenerateShortURL(fullURLByte []byte) string {
	hash := md5.New()
	hash.Write(fullURLByte)
	hashBytes := hash.Sum(nil)
	return hex.EncodeToString(hashBytes)[:8]
}
