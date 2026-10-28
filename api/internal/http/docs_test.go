package http

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	nethttp "net/http"
	"net/http/httptest"
)

func TestDocsRoutes(t *testing.T) {
	specPath := filepath.Join(t.TempDir(), "api.yaml")
	const spec = "openapi: 3.1.0\ninfo:\n  title: WallAce API\n"
	if err := os.WriteFile(specPath, []byte(spec), 0o600); err != nil {
		t.Fatalf("write spec: %v", err)
	}

	router := NewRouter(Dependencies{
		Pinger:          fakePinger{},
		Store:           &fakeStore{},
		Token:           testToken,
		Notifier:        &fakeNotifier{},
		Hub:             NewHub(),
		OpenAPISpecPath: specPath,
	})

	uiResponse := httptest.NewRecorder()
	uiRequest := httptest.NewRequest(nethttp.MethodGet, "/docs", nil)
	router.ServeHTTP(uiResponse, uiRequest)
	if uiResponse.Code != nethttp.StatusOK {
		t.Fatalf("Swagger UI status = %d, want %d", uiResponse.Code, nethttp.StatusOK)
	}
	if contentType := uiResponse.Header().Get("Content-Type"); !strings.HasPrefix(contentType, "text/html") {
		t.Fatalf("Swagger UI content type = %q, want text/html", contentType)
	}
	if !strings.Contains(uiResponse.Body.String(), `url: "/docs/openapi.yaml"`) {
		t.Fatal("Swagger UI does not reference the OpenAPI endpoint")
	}

	specResponse := httptest.NewRecorder()
	specRequest := httptest.NewRequest(nethttp.MethodGet, "/docs/openapi.yaml", nil)
	router.ServeHTTP(specResponse, specRequest)
	if specResponse.Code != nethttp.StatusOK {
		t.Fatalf("OpenAPI spec status = %d, want %d", specResponse.Code, nethttp.StatusOK)
	}
	if contentType := specResponse.Header().Get("Content-Type"); !strings.HasPrefix(contentType, "application/yaml") {
		t.Fatalf("OpenAPI spec content type = %q, want application/yaml", contentType)
	}
	if specResponse.Body.String() != spec {
		t.Fatalf("OpenAPI spec = %q, want %q", specResponse.Body.String(), spec)
	}
}
