package handler

import (
	"encoding/json"
	"io"
	"net/http"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponce struct {
	Result string `json:"result"`
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
		return
	}

	var requstPayload ShortenRequest

	if err := json.Unmarshal(body, &requstPayload); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
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
		return
	}

	res.Header().Set("content-type", "application/json")
	res.WriteHeader(http.StatusCreated)
	res.Write(response)
}
