package http

import (
	"context"
	"encoding/json"
	nethttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/caiohperlin/wallace/api/internal/database"
)

const testToken = "secret"

type fakeStore struct {
	devices        []database.Device
	config         database.DeviceConfig
	detections     []database.Detection
	listErr        error
	renameErr      error
	configErr      error
	heartbeatErr   error
	statusErr      error
	detectionsErr  error
	createErr      error
	renameCalled   bool
	heartbeatCalls int
}

func (f *fakeStore) ListDevices(context.Context) ([]database.Device, error) {
	return f.devices, f.listErr
}

func (f *fakeStore) RenameDevice(_ context.Context, id, name string) (database.Device, error) {
	f.renameCalled = true
	if f.renameErr != nil {
		return database.Device{}, f.renameErr
	}
	return database.Device{ID: id, Name: name}, nil
}

func (f *fakeStore) GetConfig(context.Context, string) (database.DeviceConfig, error) {
	return f.config, f.configErr
}

func (f *fakeStore) UpdateConfig(_ context.Context, _ string, cfg database.DeviceConfig) (database.DeviceConfig, error) {
	if f.configErr != nil {
		return database.DeviceConfig{}, f.configErr
	}
	return cfg, nil
}

func (f *fakeStore) RecordHeartbeat(context.Context, string) error {
	f.heartbeatCalls++
	return f.heartbeatErr
}

func (f *fakeStore) GetStatus(context.Context, string) (database.DeviceStatus, error) {
	return database.DeviceStatus{DeviceConfig: f.config}, f.statusErr
}

func (f *fakeStore) ListDetections(context.Context, string, int, int) ([]database.Detection, error) {
	return f.detections, f.detectionsErr
}

func (f *fakeStore) CreateDetection(_ context.Context, _ string, ts int64) (database.Detection, error) {
	if f.createErr != nil {
		return database.Detection{}, f.createErr
	}
	return database.Detection{ID: "det-1", TS: ts}, nil
}

type fakeNotifier struct {
	events []DashboardDetectionEvent
}

func (f *fakeNotifier) Broadcast(event DashboardDetectionEvent) {
	f.events = append(f.events, event)
}

func newTestRouter(store DeviceStore, notifier DetectionNotifier) nethttp.Handler {
	return NewRouter(Dependencies{
		Pinger:   fakePinger{},
		Store:    store,
		Token:    testToken,
		Notifier: notifier,
		Hub:      NewHub(),
	})
}

func doRequest(t *testing.T, router nethttp.Handler, method, path, body, token string) *httptest.ResponseRecorder {
	t.Helper()

	var reader *strings.Reader
	if body == "" {
		reader = strings.NewReader("")
	} else {
		reader = strings.NewReader(body)
	}

	req := httptest.NewRequest(method, path, reader)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func TestDeviceRoutesRequireToken(t *testing.T) {
	router := newTestRouter(&fakeStore{}, &fakeNotifier{})

	rec := doRequest(t, router, nethttp.MethodGet, "/api/devices", "", "")
	if rec.Code != nethttp.StatusUnauthorized {
		t.Fatalf("status without token = %d, want %d", rec.Code, nethttp.StatusUnauthorized)
	}
	if challenge := rec.Header().Get("WWW-Authenticate"); challenge != `Bearer realm="wallace"` {
		t.Fatalf("authentication challenge = %q, want Bearer realm=wallace", challenge)
	}

	rec = doRequest(t, router, nethttp.MethodGet, "/api/devices", "", "wrong")
	if rec.Code != nethttp.StatusUnauthorized {
		t.Fatalf("status with wrong token = %d, want %d", rec.Code, nethttp.StatusUnauthorized)
	}
}

func TestListDevices(t *testing.T) {
	store := &fakeStore{
		devices: []database.Device{{ID: "A4:CF:12:34:56:78", Name: "Front door", Online: true}},
	}
	router := newTestRouter(store, &fakeNotifier{})

	rec := doRequest(t, router, nethttp.MethodGet, "/api/devices", "", testToken)
	if rec.Code != nethttp.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, nethttp.StatusOK)
	}

	var devices []database.Device
	if err := json.NewDecoder(rec.Body).Decode(&devices); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if len(devices) != 1 || devices[0].ID != "A4:CF:12:34:56:78" {
		t.Fatalf("unexpected devices: %+v", devices)
	}
}

func TestRenameDevice(t *testing.T) {
	store := &fakeStore{}
	router := newTestRouter(store, &fakeNotifier{})

	rec := doRequest(t, router, nethttp.MethodPatch, "/api/devices/A4:CF:12:34:56:78", `{"name":"Front door"}`, testToken)
	if rec.Code != nethttp.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, nethttp.StatusOK)
	}
	if !store.renameCalled {
		t.Fatal("RenameDevice was not called")
	}
}

func TestCreateDetection(t *testing.T) {
	store := &fakeStore{}
	notifier := &fakeNotifier{}
	router := newTestRouter(store, notifier)

	rec := doRequest(t, router, nethttp.MethodPost, "/api/devices/A4:CF:12:34:56:78/detections", `{"ts":1732300000}`, testToken)
	if rec.Code != nethttp.StatusCreated {
		t.Fatalf("status = %d, want %d", rec.Code, nethttp.StatusCreated)
	}

	if len(notifier.events) != 1 {
		t.Fatalf("broadcast events = %d, want 1", len(notifier.events))
	}
	event := notifier.events[0]
	if event.Type != "detection" || event.DeviceID != "A4:CF:12:34:56:78" || event.TS != 1732300000 {
		t.Fatalf("unexpected event: %+v", event)
	}
}

func TestCreateDetectionRejectsInvalidTimestamp(t *testing.T) {
	router := newTestRouter(&fakeStore{}, &fakeNotifier{})

	rec := doRequest(t, router, nethttp.MethodPost, "/api/devices/A4:CF:12:34:56:78/detections", `{"ts":0}`, testToken)
	if rec.Code != nethttp.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, nethttp.StatusBadRequest)
	}
}

func TestListDetectionsRejectsInvalidLimit(t *testing.T) {
	router := newTestRouter(&fakeStore{}, &fakeNotifier{})

	rec := doRequest(t, router, nethttp.MethodGet, "/api/devices/A4:CF:12:34:56:78/detections?limit=0", "", testToken)
	if rec.Code != nethttp.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, nethttp.StatusBadRequest)
	}
}

func TestHeartbeat(t *testing.T) {
	store := &fakeStore{}
	router := newTestRouter(store, &fakeNotifier{})

	rec := doRequest(t, router, nethttp.MethodPost, "/api/devices/A4:CF:12:34:56:78/heartbeat", "", testToken)
	if rec.Code != nethttp.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, nethttp.StatusOK)
	}
	if store.heartbeatCalls != 1 {
		t.Fatalf("heartbeat calls = %d, want 1", store.heartbeatCalls)
	}
}
