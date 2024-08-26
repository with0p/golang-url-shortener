package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/with0p/golang-url-shortener.git/internal/auth"
	"github.com/with0p/golang-url-shortener.git/internal/logger"
)

func (handler *URLHandler) DeleteURLs(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodDelete {
		http.Error(res, "Not a Delete requests", http.StatusMethodNotAllowed)
		return
	}

	userID, err := auth.GetUserIDFromCookie(req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusUnauthorized)
		return
	}

	defer req.Body.Close()
	body, bodyReadError := io.ReadAll(req.Body)
	if bodyReadError != nil {
		http.Error(res, bodyReadError.Error(), http.StatusBadRequest)
		logger.LogError(bodyReadError)
		return
	}

	var requestPayload []string
	if err := json.Unmarshal(body, &requestPayload); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		logger.LogError(err)
		return
	}

	go handler.service.DeleteUserURLs(userID, requestPayload)

	res.Header().Set("content-type", "application/json")
	res.WriteHeader(http.StatusAccepted)
}
