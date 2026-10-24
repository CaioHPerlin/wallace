package http

import (
	"context"
	"errors"
	nethttp "net/http"
	"net/http/httptest"
	"testing"
)

type fakePinger struct {
	err error
}

func (f fakePinger) Ping(context.Context) error {
	return f.err
}

func TestHealthRoutes(t *testing.T) {
	tests := []struct {
		name       string
		pinger     Pinger
		path       string
		wantStatus int
		wantBody   string
	}{
		{name: "live", pinger: fakePinger{}, path: "/live", wantStatus: nethttp.StatusOK, wantBody: `{"status":"ok"}` + "\n"},
		{name: "ready ok", pinger: fakePinger{}, path: "/ready", wantStatus: nethttp.StatusOK, wantBody: `{"status":"ok"}` + "\n"},
		{name: "ready db down", pinger: fakePinger{err: errors.New("db unreachable")}, path: "/ready", wantStatus: nethttp.StatusServiceUnavailable, wantBody: `{"status":"unavailable"}` + "\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := NewHealthHandler(tt.pinger).Routes()
			rr := httptest.NewRecorder()
			req := httptest.NewRequest(nethttp.MethodGet, tt.path, nil)
			router.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rr.Code, tt.wantStatus)
			}
			if got := rr.Body.String(); got != tt.wantBody {
				t.Fatalf("body = %q, want %q", got, tt.wantBody)
			}
		})
	}
}
