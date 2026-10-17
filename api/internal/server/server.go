package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

const shutdownTimeout = 10 * time.Second

// Server owns the HTTP listener and its graceful shutdown lifecycle.
type Server struct {
	httpServer *http.Server
}

// New creates an HTTP server for the given address and handler.
func New(address string, handler http.Handler) *Server {
	return &Server{httpServer: &http.Server{
		Addr:              address,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}}
}

// Run serves requests until the listener fails or the context is canceled.
func (s *Server) Run(ctx context.Context) error {
	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("HTTP server listening on %s", s.httpServer.Addr)
		serverErrors <- s.httpServer.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("serve HTTP: %w", err)
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown HTTP server: %w", err)
		}
		if err := <-serverErrors; !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("serve HTTP during shutdown: %w", err)
		}
		return nil
	}
}
