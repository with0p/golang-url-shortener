package handler

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/with0p/golang-url-shortener.git/internal/config"
	"github.com/with0p/golang-url-shortener.git/internal/service"
	"github.com/with0p/golang-url-shortener.git/internal/storage"
)

func getInMemoryMocks() (*URLHandler, storage.Storage, *config.Config) {
	inMemoryStorage := storage.NewInMemoryStorage(map[string]string{})
	configuration := config.GetConfig()
	service := service.NewShortURLService(inMemoryStorage, configuration.ShortURL)
	handler := NewURLHandler(service)

	return handler, inMemoryStorage, configuration
}

func makeRequest(method string, path string, body []byte, router http.Handler) *http.Response {
	request := httptest.NewRequest(method, path, bytes.NewReader(body))
	request.Header.Set("content-type", "text/plain")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	return w.Result()
}
func TestGetTrueURL(t *testing.T) {
	type testData struct {
		method   string
		endpoint string
		shortURL string
		trueURL  string
	}
	type expectedData struct {
		status         int
		locationHeader string
	}

	tests := []struct {
		name         string
		testData     testData
		expectedData expectedData
	}{
		{
			name: "Check redirect status code",
			testData: testData{
				method:   http.MethodGet,
				endpoint: "/shorturl0",
				shortURL: "shorturl0",
				trueURL:  "http://github.com/with0p/golang-url-shortener.git",
			},
			expectedData: expectedData{
				status:         http.StatusTemporaryRedirect,
				locationHeader: "http://github.com/with0p/golang-url-shortener.git",
			},
		},
		{
			name: "Check not found status code",
			testData: testData{
				method:   http.MethodGet,
				shortURL: "",
				trueURL:  "",
				endpoint: "/shorturl0",
			},
			expectedData: expectedData{
				status:         http.StatusNotFound,
				locationHeader: "",
			},
		},
		{
			name: "Check wrong http method",
			testData: testData{
				method:   http.MethodPost,
				shortURL: "",
				trueURL:  "",
				endpoint: "/shorturl0",
			},
			expectedData: expectedData{
				status:         http.StatusMethodNotAllowed,
				locationHeader: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			URLHandler, storage, _ := getInMemoryMocks()
			router := URLHandler.GetHTTPHandler()

			if tt.testData.shortURL != "" && tt.testData.trueURL != "" {
				storage.Write(tt.testData.shortURL, tt.testData.trueURL)

			}

			res := makeRequest(tt.testData.method, tt.testData.endpoint, nil, router)
			defer res.Body.Close()

			assert.Equal(t, tt.expectedData.status, res.StatusCode)
			assert.Equal(t, tt.expectedData.locationHeader, res.Header.Get("Location"))
		})
	}

}

func TestURLShortener(t *testing.T) {
	type testData struct {
		method      string
		contentType string
		requestBody []byte
	}
	type expectedData struct {
		status      int
		shortURLId  string
		contentType string
	}

	tests := []struct {
		name         string
		testData     testData
		expectedData expectedData
	}{
		{
			name: "Check correctly created url",
			testData: testData{
				method:      http.MethodPost,
				contentType: "text/plain",
				requestBody: []byte("https://practicum.yandex.kz/"),
			},
			expectedData: expectedData{
				shortURLId:  "a0c7ecc8",
				status:      http.StatusCreated,
				contentType: "text/plain",
			},
		},
		{
			name: "Check wrong http method",
			testData: testData{
				method:      http.MethodGet,
				contentType: "text/plain",
				requestBody: []byte("https://practicum.yandex.kz/"),
			},
			expectedData: expectedData{
				shortURLId:  "",
				status:      http.StatusMethodNotAllowed,
				contentType: "text/plain",
			},
		},
		{
			name: "Check empty body",
			testData: testData{
				method:      http.MethodPost,
				contentType: "text/plain",
				requestBody: nil,
			},
			expectedData: expectedData{
				shortURLId:  "",
				status:      http.StatusBadRequest,
				contentType: "text/plain",
			},
		},
		{
			name: "Check not url body",
			testData: testData{
				method:      http.MethodPost,
				contentType: "text/plain",
				requestBody: []byte("httpspracticum.yandex.kz/"),
			},
			expectedData: expectedData{
				shortURLId:  "",
				status:      http.StatusBadRequest,
				contentType: "text/plain",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			URLHandler, _, configuration := getInMemoryMocks()
			router := URLHandler.GetHTTPHandler()

			res := makeRequest(tt.testData.method, "/", tt.testData.requestBody, router)
			defer res.Body.Close()

			body, _ := io.ReadAll(res.Body)
			assert.Equal(t, tt.expectedData.status, res.StatusCode)

			if tt.expectedData.status == http.StatusCreated {
				assert.Equal(t, configuration.ShortURL+"/"+tt.expectedData.shortURLId, string(body))
				assert.Equal(t, tt.testData.contentType, res.Header.Get("content-type"))
			}
		})
	}
}
