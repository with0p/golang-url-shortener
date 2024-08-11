package handler

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/with0p/golang-url-shortener.git/internal/common-types"
	"github.com/with0p/golang-url-shortener.git/internal/config"
	"github.com/with0p/golang-url-shortener.git/internal/mock"
	"github.com/with0p/golang-url-shortener.git/internal/service"
	"github.com/with0p/golang-url-shortener.git/internal/storage"
)

func getInMemoryMocks() *URLHandler {
	inMemoryStorage := storage.NewInMemoryStorage(map[string]string{})
	service := service.NewShortURLService(inMemoryStorage, config.MockConfiguration.ShortURL)
	handler := NewURLHandler(service)

	return handler
}

func getHandlerGetTrueURLMock(ctrl *gomock.Controller, key string, value string) *URLHandler {
	mockService := mock.NewMockService(ctrl)
	mockService.EXPECT().GetTrueURL(key).Return(value, nil)

	return NewURLHandler(mockService)
}

func getHandlerMakeShortURLMock(ctrl *gomock.Controller, key string, value string) *URLHandler {
	mockService := mock.NewMockService(ctrl)
	mockService.EXPECT().MakeShortURL(key).Return(value, nil)

	return NewURLHandler(mockService)
}

func getHandlerMakeShortURLBatchMock(ctrl *gomock.Controller, key *[]commontypes.RecordToBatch, value *[]commontypes.BatchRecord) *URLHandler {
	mockService := mock.NewMockService(ctrl)
	mockService.EXPECT().MakeShortURLBatch(key).Return(value, nil)

	return NewURLHandler(mockService)
}

func getDefaultHandler() *URLHandler {
	return getInMemoryMocks()
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
		errorExpected  bool
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
				errorExpected:  false,
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
				errorExpected:  true,
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
				errorExpected:  true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			URLHandler := getDefaultHandler()

			if !tt.expectedData.errorExpected {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				URLHandler = getHandlerGetTrueURLMock(ctrl, tt.testData.shortURL, tt.testData.trueURL)
			}

			router := URLHandler.GetHTTPHandler(nil)
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
		trueURL     string
	}
	type expectedData struct {
		status        int
		shortURL      string
		contentType   string
		errorExpected bool
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
				trueURL:     "https://practicum.yandex.kz/",
			},
			expectedData: expectedData{
				shortURL:      "http://localhost:8080/a0c7ecc8",
				status:        http.StatusCreated,
				contentType:   "text/plain",
				errorExpected: false,
			},
		},
		{
			name: "Check wrong http method",
			testData: testData{
				method:      http.MethodGet,
				contentType: "text/plain",
				trueURL:     "https://practicum.yandex.kz/",
			},
			expectedData: expectedData{
				shortURL:      "",
				status:        http.StatusMethodNotAllowed,
				contentType:   "text/plain",
				errorExpected: true,
			},
		},
		{
			name: "Check empty body",
			testData: testData{
				method:      http.MethodPost,
				contentType: "text/plain",
				trueURL:     "",
			},
			expectedData: expectedData{
				shortURL:      "",
				status:        http.StatusBadRequest,
				contentType:   "text/plain",
				errorExpected: true,
			},
		},
		{
			name: "Check not url body",
			testData: testData{
				method:      http.MethodPost,
				contentType: "text/plain",
				trueURL:     "httpspracticum.yandex.kz/",
			},
			expectedData: expectedData{
				shortURL:      "",
				status:        http.StatusBadRequest,
				contentType:   "text/plain",
				errorExpected: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			URLHandler := getDefaultHandler()

			if !tt.expectedData.errorExpected {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				URLHandler = getHandlerMakeShortURLMock(ctrl, tt.testData.trueURL, tt.expectedData.shortURL)
			}

			router := URLHandler.GetHTTPHandler(nil)
			res := makeRequest(tt.testData.method, "/", []byte(tt.testData.trueURL), tt.testData.contentType, router)
			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			require.Nil(t, err)
			assert.Equal(t, tt.expectedData.status, res.StatusCode)

			if tt.expectedData.status == http.StatusCreated {
				assert.Equal(t, tt.expectedData.shortURL, string(body))
				assert.Equal(t, tt.testData.contentType, res.Header.Get("content-type"))
			}
		})
	}
}

func TestShorten(t *testing.T) {
	type testData struct {
		method         string
		contentType    string
		requestPayload string
		trueURL        string
	}
	type expectedData struct {
		status          int
		contentType     string
		responsePayload string
		shortURL        string
		errorExpected   bool
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
				trueURL:        "https://practicum.yandex.kz/",
			},
			expectedData: expectedData{
				status:          http.StatusCreated,
				contentType:     "application/json",
				responsePayload: `{"result":"http://localhost:8080/a0c7ecc8"}`,
				shortURL:        "http://localhost:8080/a0c7ecc8",
				errorExpected:   false,
			},
		},
		{
			name: "Check wrong content type",
			testData: testData{
				method:         http.MethodPost,
				contentType:    "plain/text",
				requestPayload: `{"url":"https://practicum.yandex.kz/"}`,
				trueURL:        "https://practicum.yandex.kz/",
			},
			expectedData: expectedData{
				status:          http.StatusBadRequest,
				contentType:     "",
				responsePayload: "",
				shortURL:        "",
				errorExpected:   true,
			},
		},
		{
			name: "Check wrong payload structure",
			testData: testData{
				method:         http.MethodPost,
				contentType:    "application/json",
				requestPayload: `{"link":"https://practicum.yandex.kz/"}`,
				trueURL:        "https://practicum.yandex.kz/",
			},
			expectedData: expectedData{
				status:          http.StatusBadRequest,
				contentType:     "",
				responsePayload: "",
				shortURL:        "",
				errorExpected:   true,
			},
		},
		{
			name: "Check invalid url in payload",
			testData: testData{
				method:         http.MethodPost,
				contentType:    "application/json",
				requestPayload: `{"url":"httpracticum.yandex.kz/"}`,
				trueURL:        "https://practicum.yandex.kz/",
			},
			expectedData: expectedData{
				status:          http.StatusBadRequest,
				contentType:     "",
				responsePayload: "",
				shortURL:        "",
				errorExpected:   true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			URLHandler := getDefaultHandler()

			if !tt.expectedData.errorExpected {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				URLHandler = getHandlerMakeShortURLMock(ctrl, tt.testData.trueURL, tt.expectedData.shortURL)
			}

			router := URLHandler.GetHTTPHandler(nil)

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

func TestShortenBatch(t *testing.T) {
	type testData struct {
		method          string
		contentType     string
		requestPayload  string
		trueURLsToBatch *[]commontypes.RecordToBatch
	}
	type expectedData struct {
		status           int
		contentType      string
		responsePayload  string
		shortURLsBatched *[]commontypes.BatchRecord
		errorExpected    bool
	}

	tests := []struct {
		name         string
		testData     testData
		expectedData expectedData
	}{
		{
			name: "Check correctly shortened url batch",
			testData: testData{
				method:         http.MethodPost,
				contentType:    "application/json",
				requestPayload: `[{"correlation_id": "1","original_url": "https://practicum.yandex.fr/"},{"correlation_id": "2","original_url": "https://practicum.yandex.com/"}]`,
				trueURLsToBatch: &[]commontypes.RecordToBatch{
					{
						ID:      "1",
						FullURL: "https://practicum.yandex.fr/",
					},
					{
						ID:      "2",
						FullURL: "https://practicum.yandex.com/",
					},
				},
			},
			expectedData: expectedData{
				status:          http.StatusCreated,
				contentType:     "application/json",
				responsePayload: `[{"correlation_id":"1","short_url":"http://localhost:8080/e61c85d5"},{"correlation_id":"2","short_url":"http://localhost:8080/f17e9784"}]`,
				shortURLsBatched: &[]commontypes.BatchRecord{
					{
						ID:          "1",
						FullURL:     "https://practicum.yandex.fr/",
						ShortURL:    "http://localhost:8080/e61c85d5",
						ShortURLKey: "e61c85d5",
					},
					{
						ID:          "2",
						FullURL:     "https://practicum.yandex.com/",
						ShortURL:    "http://localhost:8080/f17e9784",
						ShortURLKey: "f17e9784",
					},
				},
				errorExpected: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			URLHandler := getDefaultHandler()

			if !tt.expectedData.errorExpected {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				URLHandler = getHandlerMakeShortURLBatchMock(ctrl, tt.testData.trueURLsToBatch, tt.expectedData.shortURLsBatched)
			}

			router := URLHandler.GetHTTPHandler(nil)

			res := makeRequest(tt.testData.method, "/api/shorten/batch", []byte(tt.testData.requestPayload), tt.testData.contentType, router)
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
