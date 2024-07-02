package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/with0p/golang-url-shortener.git/cmd/shortener/storage"
)

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

	router := ServerRouter()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage.InitMap()

			if tt.testData.shortURL != "" && tt.testData.trueURL != "" {
				storage.GetURLMap().Set(tt.testData.shortURL, tt.testData.trueURL)
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
		status       int
		responseBody string
		contentType  string
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
				responseBody: "http://localhost:8080/" + GenerateShortURL([]byte("https://practicum.yandex.kz/")),
				status:       http.StatusCreated,
				contentType:  "text/plain",
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
				responseBody: "",
				status:       http.StatusMethodNotAllowed,
				contentType:  "text/plain",
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
				responseBody: "",
				status:       http.StatusBadRequest,
				contentType:  "text/plain",
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
				responseBody: "",
				status:       http.StatusBadRequest,
				contentType:  "text/plain",
			},
		},
	}

	router := ServerRouter()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage.InitMap()

			res := makeRequest(tt.testData.method, "/", tt.testData.requestBody, router)
			defer res.Body.Close()

			body, _ := io.ReadAll(res.Body)
			assert.Equal(t, tt.expectedData.status, res.StatusCode)

			if tt.expectedData.status == http.StatusCreated {
				assert.Equal(t, tt.expectedData.responseBody, string(body))
				assert.Equal(t, tt.testData.contentType, res.Header.Get("content-type"))
			}
		})
	}

	doubleRequestTest := struct {
		name                 string
		testData             testData
		expectedStorageKey   string
		expectedStorageValue string
		expectedStatusCode   int
	}{

		name: "Check two same Posts",
		testData: testData{
			method:      http.MethodPost,
			contentType: "text/plain",
			requestBody: []byte("https://practicum.yandex.kz/"),
		},
		expectedStorageKey:   GenerateShortURL([]byte("https://practicum.yandex.kz/")),
		expectedStorageValue: "https://practicum.yandex.kz/",
		expectedStatusCode:   http.StatusCreated,
	}

	t.Run(doubleRequestTest.name, func(t *testing.T) {
		storage.InitMap()
		testStorage := storage.GetURLMap()

		response1 := makeRequest(doubleRequestTest.testData.method, "/", doubleRequestTest.testData.requestBody, router)
		defer response1.Body.Close()

		storageValueFirstRead, _ := testStorage.Get(doubleRequestTest.expectedStorageKey)

		assert.Equal(t, doubleRequestTest.expectedStatusCode, response1.StatusCode)
		assert.Equal(t, doubleRequestTest.expectedStorageValue, storageValueFirstRead)
		assert.Equal(t, testStorage.GetStorageSize(), 1)

		response2 := makeRequest(doubleRequestTest.testData.method, "/", doubleRequestTest.testData.requestBody, router)
		defer response2.Body.Close()

		storageValueSecondRead, _ := testStorage.Get(doubleRequestTest.expectedStorageKey)

		assert.Equal(t, doubleRequestTest.expectedStatusCode, response2.StatusCode)
		assert.Equal(t, doubleRequestTest.expectedStorageValue, storageValueSecondRead)
		assert.Equal(t, testStorage.GetStorageSize(), 1)
	})
}

func makeRequest(method string, path string, body []byte, router chi.Router) *http.Response {
	request := httptest.NewRequest(method, path, bytes.NewReader(body))
	request.Header.Set("content-type", "text/plain")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	return w.Result()
}
