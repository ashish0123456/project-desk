package handler

import (
	"net/http"

	"gitlab.com/Ash3941/project-management/services/api-gateway/internal/config"
	"gitlab.com/Ash3941/project-management/services/api-gateway/internal/response"
)

type SystemHandler struct {
	cfg *config.Config
}

func NewSystemHandler(cfg *config.Config) *SystemHandler {
	return &SystemHandler{cfg: cfg}
}

func (h *SystemHandler) Live(w http.ResponseWriter, r *http.Request) {
	response.Success(w, http.StatusOK, map[string]string{"status": "alive"})
}

func (h *SystemHandler) Ready(w http.ResponseWriter, r *http.Request) {
	response.Success(w, http.StatusOK, map[string]string{"status": "ready"})
}

func (h *SystemHandler) Version(w http.ResponseWriter, r *http.Request) {
	response.Success(w, http.StatusOK, map[string]string{
		"app":     h.cfg.AppName,
		"version": h.cfg.Version,
		"env":     h.cfg.AppEnv,
	})
}
