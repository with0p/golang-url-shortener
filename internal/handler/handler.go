package handler

import (
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/with0p/golang-url-shortener.git/internal/middlewares"
	"github.com/with0p/golang-url-shortener.git/internal/service"
)

type URLHandler struct {
	service service.Service
}

func NewURLHandler(currentService service.Service) *URLHandler {
	return &URLHandler{service: currentService}
}

func (handler *URLHandler) GetHTTPHandler() http.Handler {
	mux := chi.NewRouter()
	mux.Post(`/`, middlewares.UseMiddlewares(handler.DoShortURL))
	mux.Get(`/{id}`, middlewares.UseMiddlewares(handler.DoGetTrueURL))
	mux.Post(`/api/shorten`, middlewares.UseMiddlewares(handler.Shorten))

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
		return
	}

	shortURL, err := handler.service.MakeShortURL(string(body))

	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(shortURL))
}

func (handler *URLHandler) DoGetTrueURL(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Not a GET requests", http.StatusMethodNotAllowed)
		return
	}

	id := chi.URLParam(req, "id")

	trueURL, error := handler.service.GetTrueURL(id)
	if error != nil {
		http.Error(res, error.Error(), http.StatusNotFound)
		return
	}
	http.Redirect(res, req, trueURL, http.StatusTemporaryRedirect)
}
