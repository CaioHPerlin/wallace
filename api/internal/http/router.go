package http

import (
	nethttp "net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(pinger Pinger) nethttp.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Mount("/api/health", NewHealthHandler(pinger).Routes())
	return r
}
