package handler

import (
	"encoding/json"
	"net/http"

	"github.com/with0p/golang-url-shortener.git/internal/auth"
	"github.com/with0p/golang-url-shortener.git/internal/logger"
)

func (handler *URLHandler) GetUserRecords(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Not a GET requests", http.StatusMethodNotAllowed)
		return
	}

	userID, err := auth.GetUserIDFromCtx(req.Context())
	if err != nil {
		http.Error(res, err.Error(), http.StatusUnauthorized)
		return
	}

	records, error := handler.service.GetAllUserRecords(req.Context(), userID)
	if error != nil {
		http.Error(res, error.Error(), http.StatusNotFound)
		return
	}

	statusCode := http.StatusOK

	if len(records) == 0 {
		statusCode = http.StatusNoContent
	}

	response, err := json.Marshal(records)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		logger.LogError(err)
		return
	}

	res.Header().Set("content-type", "application/json")
	res.WriteHeader(statusCode)
	res.Write(response)
}
