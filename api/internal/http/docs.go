package http

import (
	"embed"
	"io"
	"log"
	nethttp "net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

//go:embed swagger.html
var swaggerUI embed.FS

type DocsHandler struct {
	openAPISpecPath string
}

func NewDocsHandler(openAPISpecPath string) *DocsHandler {
	return &DocsHandler{openAPISpecPath: openAPISpecPath}
}

func (h *DocsHandler) Register(r chi.Router) {
	r.Get("/docs", h.serveUI)
	r.Get("/docs/", h.serveUI)
	r.Get("/docs/openapi.yaml", h.serveSpec)
}

func (h *DocsHandler) serveUI(w nethttp.ResponseWriter, _ *nethttp.Request) {
	content, err := swaggerUI.ReadFile("swagger.html")
	if err != nil {
		log.Printf("read Swagger UI: %v", err)
		nethttp.Error(w, "documentation unavailable", nethttp.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(content)
}

func (h *DocsHandler) serveSpec(w nethttp.ResponseWriter, _ *nethttp.Request) {
	spec, err := os.Open(h.openAPISpecPath)
	if err != nil {
		log.Printf("open OpenAPI spec: %v", err)
		nethttp.Error(w, "documentation unavailable", nethttp.StatusInternalServerError)
		return
	}
	defer spec.Close()

	w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
	if _, err := io.Copy(w, spec); err != nil {
		log.Printf("write OpenAPI spec: %v", err)
	}
}
