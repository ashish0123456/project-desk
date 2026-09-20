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

		authProxy := proxyHandler("http://auth-service:8081")
		api.Post("/auth/register", authProxy)
		api.Post("/auth/login", authProxy)
		api.Handle("/auth", authProxy)
		api.Handle("/auth/*", authProxy)

		projectProxy := proxyHandler("http://project-service:8082")
		api.Get("/projects", projectProxy)
		api.Post("/projects", projectProxy)
		api.Get("/projects/{id}", projectProxy)
		api.Post("/projects/members", projectProxy)
		api.Handle("/projects", projectProxy)
		api.Handle("/projects/*", projectProxy)

		issueProxy := proxyHandler("http://issue-service:8083")
		api.Post("/issues", issueProxy)
		api.Handle("/issues", issueProxy)
		api.Handle("/issues/*", issueProxy)
	})

	return r
}

func proxyHandler(target string) http.HandlerFunc {
	handlerProxy := handler.NewReverseProxy(target)
	return http.HandlerFunc(handlerProxy.ServeHTTP)
}
