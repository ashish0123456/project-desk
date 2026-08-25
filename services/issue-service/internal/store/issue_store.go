package store

import (
	"context"
	"database/sql"

	"gitlab.com/Ash3941/project-management/services/issue-service/internal/models"
)

type Store struct{ db *sql.DB }

func New(db *sql.DB) *Store { return &Store{db: db} }

func (s *Store) CreateIssue(ctx context.Context, i *models.Issue) error {
    query := `INSERT INTO issues (id, project_id, title, body, status, priority, assignee, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
    _, err := s.db.ExecContext(ctx, query, i.ID, i.ProjectID, i.Title, i.Body, i.Status, i.Priority, i.Assignee, i.CreatedAt)
    return err
}
