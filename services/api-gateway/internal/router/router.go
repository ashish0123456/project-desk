package router

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.com/Ash3941/project-management/services/api-gateway/internal/handler"
	"gitlab.com/Ash3941/project-management/services/api-gateway/internal/middleware"
)

func New(system *handler.SystemHandler, log *slog.Logger) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logging(log))
	r.Use(middleware.Recovery(log))

	r.Get("/health/live", system.Live)
	r.Get("/health/ready", system.Ready)

	r.Route("/api/v1", func(api chi.Router) {
		api.Get("/version", system.Version)
		// proxy to internal services
		api.Handle("/auth/*", handler.NewReverseProxy("http://auth-service:8081"))
		api.Handle("/projects/*", handler.NewReverseProxy("http://project-service:8082"))
		api.Handle("/issues/*", handler.NewReverseProxy("http://issue-service:8083"))
	})

	return r
}
