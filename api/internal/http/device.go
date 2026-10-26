package http

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"errors"
	nethttp "net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/caiohperlin/wallace/api/internal/database"
)

type DeviceStore interface {
	ListDevices(context.Context) ([]database.Device, error)
	RenameDevice(context.Context, string, string) (database.Device, error)
	GetConfig(context.Context, string) (database.DeviceConfig, error)
	UpdateConfig(context.Context, string, database.DeviceConfig) (database.DeviceConfig, error)
	RecordHeartbeat(context.Context, string) error
	GetStatus(context.Context, string) (database.DeviceStatus, error)
	ListDetections(context.Context, string, int, int) ([]database.Detection, error)
	CreateDetection(context.Context, string, int64) (database.Detection, error)
}

type DeviceHandler struct {
	store    DeviceStore
	token    string
	notifier DetectionNotifier
}

func NewDeviceHandler(store DeviceStore, token string, notifier DetectionNotifier) *DeviceHandler {
	return &DeviceHandler{store: store, token: token, notifier: notifier}
}

func (h *DeviceHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(requireToken(h.token))

	r.Get("/", h.list)
	r.Route("/{deviceID}", func(r chi.Router) {
		r.Patch("/", h.rename)
		r.Get("/config", h.getConfig)
		r.Put("/config", h.updateConfig)
		r.Post("/heartbeat", h.heartbeat)
		r.Get("/status", h.status)
		r.Get("/detections", h.listDetections)
		r.Post("/detections", h.createDetection)
	})
	return r
}

func (h *DeviceHandler) list(w nethttp.ResponseWriter, r *nethttp.Request) {
	devices, err := h.store.ListDevices(r.Context())
	if err != nil {
		writeError(w, nethttp.StatusInternalServerError, "list devices")
		return
	}
	writeJSON(w, nethttp.StatusOK, devices)
}

func (h *DeviceHandler) rename(w nethttp.ResponseWriter, r *nethttp.Request) {
	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, nethttp.StatusBadRequest, "invalid request body")
		return
	}

	device, err := h.store.RenameDevice(r.Context(), chi.URLParam(r, "deviceID"), body.Name)
	if err != nil {
		writeError(w, nethttp.StatusInternalServerError, "rename device")
		return
	}
	writeJSON(w, nethttp.StatusOK, device)
}

func (h *DeviceHandler) getConfig(w nethttp.ResponseWriter, r *nethttp.Request) {
	cfg, err := h.store.GetConfig(r.Context(), chi.URLParam(r, "deviceID"))
	if err != nil {
		writeError(w, nethttp.StatusInternalServerError, "get device config")
		return
	}
	writeJSON(w, nethttp.StatusOK, cfg)
}

func (h *DeviceHandler) updateConfig(w nethttp.ResponseWriter, r *nethttp.Request) {
	var body database.DeviceConfig
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, nethttp.StatusBadRequest, "invalid request body")
		return
	}

	cfg, err := h.store.UpdateConfig(r.Context(), chi.URLParam(r, "deviceID"), body)
	if err != nil {
		writeError(w, nethttp.StatusInternalServerError, "update device config")
		return
	}
	writeJSON(w, nethttp.StatusOK, cfg)
}

func (h *DeviceHandler) heartbeat(w nethttp.ResponseWriter, r *nethttp.Request) {
	if err := h.store.RecordHeartbeat(r.Context(), chi.URLParam(r, "deviceID")); err != nil {
		writeError(w, nethttp.StatusInternalServerError, "record heartbeat")
		return
	}
	w.WriteHeader(nethttp.StatusOK)
}

func (h *DeviceHandler) status(w nethttp.ResponseWriter, r *nethttp.Request) {
	status, err := h.store.GetStatus(r.Context(), chi.URLParam(r, "deviceID"))
	if err != nil {
		writeError(w, nethttp.StatusInternalServerError, "get device status")
		return
	}
	writeJSON(w, nethttp.StatusOK, status)
}

func (h *DeviceHandler) listDetections(w nethttp.ResponseWriter, r *nethttp.Request) {
	limit, offset, err := pagination(r)
	if err != nil {
		writeError(w, nethttp.StatusBadRequest, "invalid pagination")
		return
	}

	detections, err := h.store.ListDetections(r.Context(), chi.URLParam(r, "deviceID"), limit, offset)
	if err != nil {
		writeError(w, nethttp.StatusInternalServerError, "list detections")
		return
	}
	writeJSON(w, nethttp.StatusOK, detections)
}

func (h *DeviceHandler) createDetection(w nethttp.ResponseWriter, r *nethttp.Request) {
	var body struct {
		TS int64 `json:"ts"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, nethttp.StatusBadRequest, "invalid request body")
		return
	}
	if body.TS <= 0 {
		writeError(w, nethttp.StatusBadRequest, "ts must be a positive unix timestamp")
		return
	}

	deviceID := chi.URLParam(r, "deviceID")
	detection, err := h.store.CreateDetection(r.Context(), deviceID, body.TS)
	if err != nil {
		writeError(w, nethttp.StatusInternalServerError, "record detection")
		return
	}

	h.notifier.Broadcast(DashboardDetectionEvent{Type: "detection", DeviceID: deviceID, TS: body.TS})
	writeJSON(w, nethttp.StatusCreated, detection)
}

func pagination(r *nethttp.Request) (int, int, error) {
	limit, offset := 50, 0

	if v := r.URL.Query().Get("limit"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 {
			return 0, 0, errors.New("invalid limit")
		}
		limit = n
	}
	if v := r.URL.Query().Get("offset"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 0 {
			return 0, 0, errors.New("invalid offset")
		}
		offset = n
	}

	return limit, offset, nil
}

func requireToken(token string) func(nethttp.Handler) nethttp.Handler {
	return func(next nethttp.Handler) nethttp.Handler {
		return nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
			if !validToken(r, token) {
				writeError(w, nethttp.StatusUnauthorized, "unauthorized")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func validToken(r *nethttp.Request, token string) bool {
	const prefix = "Bearer "
	auth := r.Header.Get("Authorization")
	if len(auth) <= len(prefix) || !strings.EqualFold(auth[:len(prefix)], prefix) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(auth[len(prefix):]), []byte(token)) == 1
}
