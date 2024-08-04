package handler

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/with0p/golang-url-shortener.git/internal/config"
	"github.com/with0p/golang-url-shortener.git/internal/logger"
	"github.com/with0p/golang-url-shortener.git/internal/service"
	"github.com/with0p/golang-url-shortener.git/internal/storage"
	localFileStorage "github.com/with0p/golang-url-shortener.git/internal/storage/local-file"
)

func getConfiguration() *config.Config {
	return config.GetConfig()
}

func getInMemoryMocks() (*URLHandler, storage.Storage) {
	inMemoryStorage := storage.NewInMemoryStorage(map[string]string{})
	configuration := getConfiguration()
	service := service.NewShortURLService(inMemoryStorage, configuration.ShortURL)
	handler := NewURLHandler(service)

	return handler, inMemoryStorage
}

func getInLocalFileMocks() (*URLHandler, storage.Storage) {
	configuration := getConfiguration()
	localFileStorage, _ := localFileStorage.NewLocalFileStorage(configuration.FileStoragePath)
	service := service.NewShortURLService(localFileStorage, configuration.ShortURL)
	handler := NewURLHandler(service)

	return handler, localFileStorage
}

func getMocks() (*URLHandler, storage.Storage) {
	return getInMemoryMocks()
}

var db *sql.DB

func getDB() *sql.DB {
	if db == nil {
		dbAddress := getConfiguration().DataBaseAddress

		dataBase, dbErr := sql.Open("pgx", dbAddress)
		if dbErr != nil {
			logger.LogError(dbErr)
		}
		defer dataBase.Close()
		db = dataBase
	}

	return db
}

func makeRequest(method string, path string, body []byte, contentType string, router http.Handler) *http.Response {
	request := httptest.NewRequest(method, path, bytes.NewReader(body))
	request.Header.Set("content-type", contentType)
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
			URLHandler, storage := getMocks()

			router := URLHandler.GetHTTPHandler(getDB())

			if tt.testData.shortURL != "" && tt.testData.trueURL != "" {
				storage.Write(tt.testData.shortURL, tt.testData.trueURL)

			}

			res := makeRequest(tt.testData.method, tt.testData.endpoint, nil, "text/plain", router)
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
			URLHandler, _ := getMocks()
			configuration := getConfiguration()

			router := URLHandler.GetHTTPHandler(getDB())

			res := makeRequest(tt.testData.method, "/", tt.testData.requestBody, tt.testData.contentType, router)
			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			require.Nil(t, err)
			assert.Equal(t, tt.expectedData.status, res.StatusCode)

			if tt.expectedData.status == http.StatusCreated {
				assert.Equal(t, configuration.ShortURL+"/"+tt.expectedData.shortURLId, string(body))
				assert.Equal(t, tt.testData.contentType, res.Header.Get("content-type"))
			}
		})
	}
}

func TestShorten(t *testing.T) {
	configuration := getConfiguration()

	type testData struct {
		method         string
		contentType    string
		requestPayload string
	}
	type expectedData struct {
		status          int
		contentType     string
		responsePayload string
	}

	tests := []struct {
		name         string
		testData     testData
		expectedData expectedData
	}{
		{
			name: "Check correctly shortened url",
			testData: testData{
				method:         http.MethodPost,
				contentType:    "application/json",
				requestPayload: `{"url":"https://practicum.yandex.kz/"}`,
			},
			expectedData: expectedData{
				status:          http.StatusCreated,
				contentType:     "application/json",
				responsePayload: fmt.Sprintf(`{"result":"%s/a0c7ecc8"}`, configuration.ShortURL),
			},
		},
		{
			name: "Check wrong content type",
			testData: testData{
				method:         http.MethodPost,
				contentType:    "plain/text",
				requestPayload: `{"url":"https://practicum.yandex.kz/"}`,
			},
			expectedData: expectedData{
				status:          http.StatusBadRequest,
				contentType:     "",
				responsePayload: "",
			},
		},
		{
			name: "Check wrong payload structure",
			testData: testData{
				method:         http.MethodPost,
				contentType:    "application/json",
				requestPayload: `{"link":"https://practicum.yandex.kz/"}`,
			},
			expectedData: expectedData{
				status:          http.StatusBadRequest,
				contentType:     "",
				responsePayload: "",
			},
		},
		{
			name: "Check invalid url in payload",
			testData: testData{
				method:         http.MethodPost,
				contentType:    "application/json",
				requestPayload: `{"link":"httpracticum.yandex.kz/"}`,
			},
			expectedData: expectedData{
				status:          http.StatusBadRequest,
				contentType:     "",
				responsePayload: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			URLHandler, _ := getMocks()

			router := URLHandler.GetHTTPHandler(getDB())

			res := makeRequest(tt.testData.method, "/api/shorten", []byte(tt.testData.requestPayload), tt.testData.contentType, router)
			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			require.Nil(t, err)
			assert.Equal(t, tt.expectedData.status, res.StatusCode)

			if tt.expectedData.status == http.StatusCreated {
				assert.Equal(t, tt.expectedData.responsePayload, string(body))
				assert.Equal(t, tt.testData.contentType, res.Header.Get("content-type"))
			}
		})
	}
}
