package store

import (
	"context"
	"database/sql"
	"errors"

	"gitlab.com/Ash3941/project-management/services/auth-service/internal/models"
)

// ErrNotFound is returned when a user cannot be located in the store.
var ErrNotFound = errors.New("user not found")

// Store is a thin wrapper around *sql.DB providing user persistence.
type Store struct {
    db *sql.DB
}

// New creates a new user store.
func New(db *sql.DB) *Store { return &Store{db: db} }

// CreateUser inserts a user record. Caller should populate CreatedAt.
func (s *Store) CreateUser(ctx context.Context, u *models.User) error {
    query := `INSERT INTO users (id, email, password_hash, created_at) VALUES ($1, $2, $3, $4)`
    _, err := s.db.ExecContext(ctx, query, u.ID, u.Email, u.PasswordHash, u.CreatedAt)
    return err
}

// GetByEmail fetches a user by email. Returns ErrNotFound when missing.
func (s *Store) GetByEmail(ctx context.Context, email string) (*models.User, error) {
    var u models.User
    query := `SELECT id, email, password_hash, created_at FROM users WHERE email = $1`
    row := s.db.QueryRowContext(ctx, query, email)
    if err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrNotFound
        }
        return nil, err
    }
    return &u, nil
}
