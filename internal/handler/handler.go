package handler

import (
	"database/sql"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	customerrors "github.com/with0p/golang-url-shortener.git/internal/custom-errors"
	"github.com/with0p/golang-url-shortener.git/internal/logger"
	"github.com/with0p/golang-url-shortener.git/internal/middlewares"
	"github.com/with0p/golang-url-shortener.git/internal/service"
)

type URLHandler struct {
	service service.Service
}

func NewURLHandler(currentService service.Service) *URLHandler {
	return &URLHandler{service: currentService}
}

func (handler *URLHandler) GetHTTPHandler(db *sql.DB) http.Handler {
	mux := chi.NewRouter()
	mux.Post(`/`, middlewares.UseMiddlewares(handler.DoShortURL))
	mux.Get(`/{id}`, middlewares.UseMiddlewares(handler.DoGetTrueURL))
	mux.Post(`/api/shorten`, middlewares.UseMiddlewares(handler.Shorten))
	mux.Post(`/api/shorten/batch`, middlewares.UseMiddlewares(handler.ShortenBatch))
	mux.Get(`/ping`, getPingDB(db))

	return mux
}

func (handler *URLHandler) DoShortURL(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Not a POST requests", http.StatusMethodNotAllowed)
		return
	}

	defer req.Body.Close()
	body, bodyReadError := io.ReadAll(req.Body)
	if bodyReadError != nil {
		http.Error(res, bodyReadError.Error(), http.StatusBadRequest)
		logger.LogError(bodyReadError)
		return
	}

	statusCode := http.StatusCreated

	shortURL, serviceErr := handler.service.MakeShortURL(req.Context(), string(body))

	if serviceErr != nil {
		if errors.Is(serviceErr, customerrors.ErrUniqueKeyConstrantViolation) {
			statusCode = http.StatusConflict
		} else {
			http.Error(res, serviceErr.Error(), http.StatusBadRequest)
			return
		}
	}

	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(statusCode)
	res.Write([]byte(shortURL))
}

func (handler *URLHandler) DoGetTrueURL(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Not a GET requests", http.StatusMethodNotAllowed)
		return
	}

	id := chi.URLParam(req, "id")

	trueURL, error := handler.service.GetTrueURL(req.Context(), id)
	if error != nil {
		http.Error(res, error.Error(), http.StatusNotFound)
		return
	}
	http.Redirect(res, req, trueURL, http.StatusTemporaryRedirect)
}

func getPingDB(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if errCtx := db.PingContext(r.Context()); errCtx != nil {
			logger.LogError(errCtx)
			http.Error(w, errCtx.Error(), http.StatusInternalServerError)
			return
		}
		logger.LogInfo("DB connected")
		w.Write([]byte("DB connected"))
	}
}
