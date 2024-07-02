package handlers

import (
	"github.com/go-chi/chi/v5"
)

func ServerRouter() chi.Router {
	mux := chi.NewRouter()
	mux.Post(`/`, URLShortener)
	mux.Get(`/{id}`, GetTrueURL)
	return mux
}
