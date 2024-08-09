package handler

import (
	"encoding/json"
	"io"
	"net/http"

	commontypes "github.com/with0p/golang-url-shortener.git/internal/common-types"
	"github.com/with0p/golang-url-shortener.git/internal/logger"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponce struct {
	Result string `json:"result"`
}

type ShortenBatchRequestRecord struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type ShortenBatchResponceRecord struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func (handler *URLHandler) Shorten(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Not a POST requests", http.StatusMethodNotAllowed)
		return
	}

	if req.Header.Get("content-type") != "application/json" {
		http.Error(res, "Not a \"application/json\" content-type", http.StatusBadRequest)
		return
	}

	defer req.Body.Close()
	body, bodyReadError := io.ReadAll(req.Body)
	if bodyReadError != nil {
		http.Error(res, bodyReadError.Error(), http.StatusBadRequest)
		logger.LogError(bodyReadError)
		return
	}

	var requstPayload ShortenRequest

	if err := json.Unmarshal(body, &requstPayload); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		logger.LogError(err)
		return
	}

	shortURL, err := handler.service.MakeShortURL(requstPayload.URL)

	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	responsePayload := ShortenResponce{
		Result: shortURL,
	}

	response, err := json.Marshal(responsePayload)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		logger.LogError(err)
		return
	}

	res.Header().Set("content-type", "application/json")
	res.WriteHeader(http.StatusCreated)
	res.Write(response)
}

func (handler *URLHandler) ShortenBatch(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Not a POST requests", http.StatusMethodNotAllowed)
		return
	}

	if req.Header.Get("content-type") != "application/json" {
		http.Error(res, "Not a \"application/json\" content-type", http.StatusBadRequest)
		return
	}

	defer req.Body.Close()
	body, bodyReadError := io.ReadAll(req.Body)
	if bodyReadError != nil {
		http.Error(res, bodyReadError.Error(), http.StatusBadRequest)
		logger.LogError(bodyReadError)
		return
	}

	var requestPayload []ShortenBatchRequestRecord
	if err := json.Unmarshal(body, &requestPayload); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		logger.LogError(err)
		return
	}

	var dataToBatch []commontypes.RecordToBatch
	for _, r := range requestPayload {
		dataToBatch = append(dataToBatch, commontypes.RecordToBatch{
			ID:      r.CorrelationID,
			FullURL: r.OriginalURL,
		})
	}

	responsePayloadData, batchError := handler.service.MakeShortURLBatch(&dataToBatch)
	if batchError != nil {
		http.Error(res, batchError.Error(), http.StatusBadRequest)
		logger.LogError(batchError)
		return
	}

	var responsePayload []ShortenBatchResponceRecord
	for _, r := range *responsePayloadData {
		responsePayload = append(responsePayload, ShortenBatchResponceRecord{
			CorrelationID: r.ID,
			ShortURL:      r.ShortURL,
		})
	}

	response, err := json.Marshal(responsePayload)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		logger.LogError(err)
		return
	}

	res.Header().Set("content-type", "application/json")
	res.WriteHeader(http.StatusCreated)
	res.Write(response)
}
