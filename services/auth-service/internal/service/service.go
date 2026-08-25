package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"gitlab.com/Ash3941/project-management/services/auth-service/internal/models"
	"gitlab.com/Ash3941/project-management/services/auth-service/internal/store"
)

// Service contains business logic for authentication and user management.
type Service struct {
    store *store.Store
}

// New returns a new auth service.
func New(s *store.Store) *Service { return &Service{store: s} }

// Register creates a new user with a bcrypt-hashed password.
func (s *Service) Register(ctx context.Context, email, password string) (*models.User, error) {
    id := uuid.New().String()
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }
    u := &models.User{
        ID:           id,
        Email:        email,
        PasswordHash: string(hash),
        CreatedAt:    time.Now().UTC(),
    }
    if err := s.store.CreateUser(ctx, u); err != nil {
        return nil, err
    }
    // zero the hash before returning the user object
    u.PasswordHash = ""
    return u, nil
}

// Authenticate verifies credentials and returns the user when successful.
func (s *Service) Authenticate(ctx context.Context, email, password string) (*models.User, error) {
    u, err := s.store.GetByEmail(ctx, email)
    if err != nil {
        return nil, err
    }
    if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
        return nil, err
    }
    u.PasswordHash = ""
    return u, nil
}
