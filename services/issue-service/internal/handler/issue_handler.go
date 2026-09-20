package handler

import (
	"encoding/json"
	"net/http"

	"gitlab.com/Ash3941/project-management/services/issue-service/internal/service"
)

type Handler struct{ svc *service.Service }

func New(svc *service.Service) *Handler { return &Handler{svc: svc} }

type createReq struct {
    ProjectID string `json:"project_id"`
    Title     string `json:"title"`
    Body      string `json:"body"`
}

func (h *Handler) CreateIssue(w http.ResponseWriter, r *http.Request) {
    var req createReq
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid", http.StatusBadRequest)
        return
    }
    i, err := h.svc.CreateIssue(r.Context(), req.ProjectID, req.Title, req.Body)
    if err != nil {
        http.Error(w, "error", http.StatusInternalServerError)
        return
    }
    _ = json.NewEncoder(w).Encode(i)
}

func (h *Handler) Routes() http.Handler {
    mux := http.NewServeMux()
    mux.HandleFunc("POST /issues", h.CreateIssue)
    return mux
}
