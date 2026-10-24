package http

import (
	"context"
	"encoding/json"
	"log"
	nethttp "net/http"
	"time"
)

const readinessTimeout = 2 * time.Second

// ReadinessProbe must return nil only while the database is reachable and honor the caller's context.
type ReadinessProbe interface {
	Ping(context.Context) error
}

type RouteRegistrar interface {
	Get(pattern string, handler nethttp.HandlerFunc)
}

type healthResponse struct {
	Status string `json:"status"`
}

type HealthHandler struct {
	readiness ReadinessProbe
}

func NewHealthHandler(readiness ReadinessProbe) *HealthHandler {
	return &HealthHandler{readiness: readiness}
}

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
