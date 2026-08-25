package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"gitlab.com/Ash3941/project-management/services/project-service/internal/models"
	"gitlab.com/Ash3941/project-management/services/project-service/internal/store"
)

// Service contains project business logic.
type Service struct{ store *store.Store }

// New returns a service instance wired to the provided store.
func New(s *store.Store) *Service { return &Service{store: s} }

// CreateProject creates a new project record.
func (s *Service) CreateProject(ctx context.Context, name string) (*models.Project, error) {
    p := &models.Project{ID: uuid.New().String(), Name: name, CreatedAt: time.Now().UTC()}
    if err := s.store.CreateProject(ctx, p); err != nil {
        return nil, err
    }
    return p, nil
}

// AddMember adds a user as a member to the project with a role.
func (s *Service) AddMember(ctx context.Context, projectID, userID, role string) (*models.Member, error) {
    m := &models.Member{ID: uuid.New().String(), ProjectID: projectID, UserID: userID, Role: role, CreatedAt: time.Now().UTC()}
    if err := s.store.AddMember(ctx, m); err != nil {
        return nil, err
    }
    return m, nil
}
