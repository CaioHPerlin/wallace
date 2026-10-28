//go:build integration

package integration

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

func stackURL() string {
	if baseURL := os.Getenv("WALLACE_BASE_URL"); baseURL != "" {
		return strings.TrimRight(baseURL, "/")
	}
	return "http://localhost:5173"
}

func TestGatewayRoutes(t *testing.T) {
	client := &http.Client{Timeout: 10 * time.Second}
	for _, tt := range []struct {
		path        string
		status      int
		contentType string
		body        string
	}{
		{"/", http.StatusOK, "text/html", `id="root"`},
		{"/dashboard/test", http.StatusOK, "text/html", `id="root"`},
		{"/docs", http.StatusOK, "text/html", "WallAce · API Explorer"},
		{"/docs/", http.StatusOK, "text/html", "WallAce · API Explorer"},
		{"/docs/openapi.yaml", http.StatusOK, "application/yaml", "openapi: 3.1.0"},
		{"/api/health/live", http.StatusOK, "application/json", `"status":"ok"`},
		{"/api/devices", http.StatusUnauthorized, "application/json", `"error":"unauthorized"`},
	} {
		t.Run(tt.path, func(t *testing.T) {
			response, err := client.Get(stackURL() + tt.path)
			if err != nil {
				t.Fatal(err)
			}
			defer response.Body.Close()
			body, err := io.ReadAll(response.Body)
			if err != nil {
				t.Fatal(err)
			}
			if response.StatusCode != tt.status {
				t.Errorf("status = %d, want %d", response.StatusCode, tt.status)
			}
			if got := response.Header.Get("Content-Type"); !strings.HasPrefix(got, tt.contentType) {
				t.Errorf("content type = %q, want %q", got, tt.contentType)
			}
			if !strings.Contains(string(body), tt.body) {
				t.Errorf("response does not contain %q", tt.body)
			}
		})
	}
}

func TestGatewayAuthorization(t *testing.T) {
	token := os.Getenv("API_TOKEN")
	if token == "" {
		t.Skip("set API_TOKEN to verify authenticated requests against the running stack")
	}
	request, err := http.NewRequest(http.MethodGet, stackURL()+"/api/devices", nil)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{Timeout: 10 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("authenticated status = %d, want 200", response.StatusCode)
	}
	var devices []json.RawMessage
	if err := json.NewDecoder(response.Body).Decode(&devices); err != nil {
		t.Fatalf("decode device list: %v", err)
	}
}
