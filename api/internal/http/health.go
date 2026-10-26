package http

import (
	"context"
	nethttp "net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

const readinessTimeout = 2 * time.Second

// Pinger reports whether its backing service is reachable.
type Pinger interface {
	Ping(context.Context) error
}

type healthResponse struct {
	Status string `json:"status"`
}

type HealthHandler struct {
	pinger Pinger
}

func NewHealthHandler(pinger Pinger) *HealthHandler {
	return &HealthHandler{pinger: pinger}
}

func (h *HealthHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/live", h.live)
	r.Get("/ready", h.ready)
	return r
}

func (h *HealthHandler) live(w nethttp.ResponseWriter, _ *nethttp.Request) {
	writeJSON(w, nethttp.StatusOK, healthResponse{Status: "ok"})
}

func (h *HealthHandler) ready(w nethttp.ResponseWriter, r *nethttp.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), readinessTimeout)
	defer cancel()

	if err := h.pinger.Ping(ctx); err != nil {
		writeJSON(w, nethttp.StatusServiceUnavailable, healthResponse{Status: "unavailable"})
		return
	}
	writeJSON(w, nethttp.StatusOK, healthResponse{Status: "ok"})
}
