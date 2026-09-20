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

func (s *Store) ListProjects(ctx context.Context) ([]*models.Project, error) {
	query := `SELECT id, name, created_at FROM projects ORDER BY created_at DESC`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := make([]*models.Project, 0)
	for rows.Next() {
		p := &models.Project{}
		if err := rows.Scan(&p.ID, &p.Name, &p.CreatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return projects, nil
}

func (s *Store) GetByID(ctx context.Context, id string) (*models.Project, error) {
	query := `SELECT id, name, created_at FROM projects WHERE id = $1`
	p := &models.Project{}
	if err := s.db.QueryRowContext(ctx, query, id).Scan(&p.ID, &p.Name, &p.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return p, nil
}

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
