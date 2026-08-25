package handler

import (
	"encoding/json"
	"net/http"

	"gitlab.com/Ash3941/project-management/services/project-service/internal/service"
)

type Handler struct{ svc *service.Service }

func New(svc *service.Service) *Handler { return &Handler{svc: svc} }

type createReq struct{ Name string `json:"name"` }

func (h *Handler) CreateProject(w http.ResponseWriter, r *http.Request) {
    var req createReq
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid", http.StatusBadRequest)
        return
    }
    p, err := h.svc.CreateProject(r.Context(), req.Name)
    if err != nil {
        http.Error(w, "error", http.StatusInternalServerError)
        return
    }
    _ = json.NewEncoder(w).Encode(p)
}

type addMemberReq struct {
    ProjectID string `json:"project_id,omitempty"`
    UserID    string `json:"user_id,omitempty"`
    Role      string `json:"role,omitempty"`
}

func (h *Handler) AddMember(w http.ResponseWriter, r *http.Request) {
    var req addMemberReq
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid", http.StatusBadRequest)
        return
    }
    m, err := h.svc.AddMember(r.Context(), req.ProjectID, req.UserID, req.Role)
    if err != nil {
        http.Error(w, "error", http.StatusInternalServerError)
        return
    }
    _ = json.NewEncoder(w).Encode(m)
}

func (h *Handler) Routes() http.Handler {
    mux := http.NewServeMux()
    mux.HandleFunc("/projects", h.CreateProject)
    mux.HandleFunc("/projects/members", h.AddMember)
    return mux
}
