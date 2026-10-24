package http

import (
	"context"
	"encoding/json"
	"log"
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
	writeHealth(w, nethttp.StatusOK, "ok")
}

func (h *HealthHandler) ready(w nethttp.ResponseWriter, r *nethttp.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), readinessTimeout)
	defer cancel()

	if err := h.pinger.Ping(ctx); err != nil {
		writeHealth(w, nethttp.StatusServiceUnavailable, "unavailable")
		return
	}
	writeHealth(w, nethttp.StatusOK, "ok")
}

func writeHealth(w nethttp.ResponseWriter, statusCode int, status string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(healthResponse{Status: status}); err != nil {
		log.Printf("write health response: %v", err)
	}
}
