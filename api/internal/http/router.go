package http

import (
	nethttp "net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Dependencies struct {
	Pinger   Pinger
	Store    DeviceStore
	Token    string
	Notifier DetectionNotifier
	Hub      *Hub
}

func NewRouter(deps Dependencies) nethttp.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Mount("/api/health", NewHealthHandler(deps.Pinger).Routes())
	r.Mount("/api/devices", NewDeviceHandler(deps.Store, deps.Token, deps.Notifier).Routes())
	r.Get("/ws/dashboard", deps.Hub.ServeWS)
	return r
}
