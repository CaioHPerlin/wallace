package http

import (
	"context"
	"encoding/json"
	"log"
	nethttp "net/http"
	"time"
)

const readinessTimeout = 2 * time.Second

// ReadinessProbe is the dependency required to report API readiness.
type ReadinessProbe interface {
	Ping(context.Context) error
}

// RouteRegistrar registers GET handlers without exposing the router implementation.
type RouteRegistrar interface {
	Get(pattern string, handler nethttp.HandlerFunc)
}

type healthResponse struct {
	Status string `json:"status"`
}

// HealthHandler serves liveness and readiness checks for the API.
type HealthHandler struct {
	readiness ReadinessProbe
}

// NewHealthHandler creates health endpoints backed by the provided readiness probe.
func NewHealthHandler(readiness ReadinessProbe) *HealthHandler {
	return &HealthHandler{readiness: readiness}
}

// Register attaches health endpoints to a route registrar.
func (h *HealthHandler) Register(router RouteRegistrar) {
	router.Get("/api/health/live", h.live)
	router.Get("/api/health/ready", h.ready)
}

func (h *HealthHandler) live(w nethttp.ResponseWriter, _ *nethttp.Request) {
	writeHealth(w, nethttp.StatusOK, "ok")
}

func (h *HealthHandler) ready(w nethttp.ResponseWriter, r *nethttp.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), readinessTimeout)
	defer cancel()

	if err := h.readiness.Ping(ctx); err != nil {
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
