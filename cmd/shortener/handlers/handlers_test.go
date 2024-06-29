package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/with0p/golang-url-shortener.git/cmd/shortener/storage"
)

func TestGetTrueURL(t *testing.T) {
	type args struct {
		endpoint string
		status   int
		key      string
		value    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Check redirect status code",
			args: args{
				endpoint: "/shorturl0",
				status:   http.StatusTemporaryRedirect,
				key:      "shorturl0",
				value:    "http://github.com/with0p/golang-url-shortener.git",
			},
		},
		{
			name: "Check not found status code",
			args: args{
				endpoint: "/shorturl01",
				status:   http.StatusNotFound,
				key:      "",
				value:    "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage.InitMap()

			if tt.args.key != "" && tt.args.value != "" {
				storage.GetURLMap().Set(tt.args.key, tt.args.value)
			}

			request := httptest.NewRequest(http.MethodGet, tt.args.endpoint, nil)
			w := httptest.NewRecorder()
			GetTrueURL(w, request)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.args.status, res.StatusCode)

			// if res.StatusCode == 307 {
			// 	id := strings.Split(tt.args.endpoint, "/")[1]
			// 	trueURL, _ := storage.GetURLMap().Get(id)
			// 	assert.Equal(t, trueURL, res.Header.Get("Location"))
			// }
		})
	}
}

func TestURLShortener(t *testing.T) {
	type args struct {
		requestBody    string
		responseBody   string
		responseStatus int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Check redirect status code",
			args: args{
				requestBody:    "https://practicum.yandex.kz/",
				responseBody:   "http://localhost:8080/" + GenerateShortURL([]byte("https://practicum.yandex.kz/")),
				responseStatus: http.StatusCreated,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage.InitMap()
			request := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(tt.args.requestBody)))
			w := httptest.NewRecorder()
			URLShortener(w, request)

			res := w.Result()

			defer res.Body.Close()
			body, _ := io.ReadAll(res.Body)

			assert.Equal(t, tt.args.responseBody, string(body))
		})
	}
}
