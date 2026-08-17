package router

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"gitlab.com/Ash3941/go-pm-platform/internal/handler"
	"gitlab.com/Ash3941/go-pm-platform/internal/middleware"
)

func New(system *handler.SystemHandler, log *slog.Logger) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recovery(log))

	r.Get("/health", system.HealthCheck)

	r.Route("api/v1", func(api chi.Router) {
		api.Get("/version", system.Version)
	})

	return r
}