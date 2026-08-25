package server

import (
	"context"
	"net/http"
)

// Server wraps the underlying HTTP server used by each service.
type Server struct {
	httpServer *http.Server
}

// New constructs an HTTP server with the supplied handler and port.
func New(handler http.Handler, port string) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    ":" + port,
			Handler: handler,
		},
	}
}

// Start begins serving HTTP traffic.
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully terminates the HTTP server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
