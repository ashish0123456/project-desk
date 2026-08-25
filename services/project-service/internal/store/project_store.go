package store

import (
	"context"
	"database/sql"
	"errors"

	"gitlab.com/Ash3941/project-management/services/project-service/internal/models"
)

var ErrNotFound = errors.New("project not found")

type Store struct{ db *sql.DB }

func New(db *sql.DB) *Store { return &Store{db: db} }

func (s *Store) CreateProject(ctx context.Context, p *models.Project) error {
    query := `INSERT INTO projects (id, name, created_at) VALUES ($1, $2, $3)`
    _, err := s.db.ExecContext(ctx, query, p.ID, p.Name, p.CreatedAt)
    return err
}

func (s *Store) AddMember(ctx context.Context, m *models.Member) error {
    query := `INSERT INTO project_members (id, project_id, user_id, role, created_at) VALUES ($1, $2, $3, $4, $5)`
    _, err := s.db.ExecContext(ctx, query, m.ID, m.ProjectID, m.UserID, m.Role, m.CreatedAt)
    return err
}
