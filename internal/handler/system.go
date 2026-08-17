package handler

import (
	"database/sql"
	"net/http"

	"gitlab.com/Ash3941/go-pm-platform/internal/config"
	"gitlab.com/Ash3941/go-pm-platform/internal/response"
)

type SystemHandler struct {
	cfg *config.Config
	db *sql.DB
}


func NewSystemHandler(cfg *config.Config, db *sql.DB) *SystemHandler {
	return &SystemHandler{
		cfg: cfg,
		db: db,
	}
}

func (h *SystemHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if err := h.db.PingContext(r.Context()); err != nil {
		response.Error(w, http.StatusServiceUnavailable, "database unavailable")
		return
	}
	response.Success(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *SystemHandler) Version(w http.ResponseWriter, r *http.Request) {
	response.Success(w, http.StatusOK, map[string]string{
		"app": h.cfg.AppName,
		"version": h.cfg.Version,
	})
}