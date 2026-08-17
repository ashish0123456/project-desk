package middleware

import (
	"log/slog"
	"net/http"

	"gitlab.com/Ash3941/go-pm-platform/internal/response"
)	

func Recovery(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			defer func() {

				if err := recover(); err != nil {

					log.Error("panic recovered",
						"error", err,
						"path", r.URL.Path,
					)

					response.Error(
						w,
						http.StatusInternalServerError,
						"internal server error",
					)
				}

			}()

			next.ServeHTTP(w, r)

		})
	}
}