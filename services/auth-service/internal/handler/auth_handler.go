package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"gitlab.com/Ash3941/project-management/services/auth-service/internal/jwt"
	"gitlab.com/Ash3941/project-management/services/auth-service/internal/service"
)

type Handler struct {
    svc *service.Service
    jwt *jwt.Manager
}

func New(svc *service.Service, jwtm *jwt.Manager) *Handler { return &Handler{svc: svc, jwt: jwtm} }

type registerRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type registerResponse struct {
    ID    string `json:"id"`
    Email string `json:"email"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
    var req registerRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request", http.StatusBadRequest)
        return
    }
    ctx := r.Context()
    u, err := h.svc.Register(ctx, req.Email, req.Password)
    if err != nil {
        // hide internal errors from clients for security
        http.Error(w, "unable to register", http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusCreated)
    _ = json.NewEncoder(w).Encode(registerResponse{ID: u.ID, Email: u.Email})
}

type loginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type loginResponse struct {
    AccessToken  string    `json:"access_token"`
    RefreshToken string    `json:"refresh_token"`
    ExpiresAt    time.Time `json:"expires_at"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
    var req loginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request", http.StatusBadRequest)
        return
    }
    ctx := r.Context()
    u, err := h.svc.Authenticate(ctx, req.Email, req.Password)
    if err != nil {
        // avoid leaking whether email existed
        http.Error(w, "invalid credentials", http.StatusUnauthorized)
        return
    }
    access, err := h.jwt.NewAccessToken(u.ID, 15*time.Minute)
    if err != nil {
        http.Error(w, "token error", http.StatusInternalServerError)
        return
    }
    refresh, err := h.jwt.NewRefreshToken(u.ID, 7*24*time.Hour)
    if err != nil {
        http.Error(w, "token error", http.StatusInternalServerError)
        return
    }
    resp := loginResponse{AccessToken: access, RefreshToken: refresh, ExpiresAt: time.Now().Add(15 * time.Minute)}
    _ = json.NewEncoder(w).Encode(resp)
}

func (h *Handler) Routes() http.Handler {
    mux := http.NewServeMux()
    mux.HandleFunc("/register", h.Register)
    mux.HandleFunc("/login", h.Login)
    return mux
}
