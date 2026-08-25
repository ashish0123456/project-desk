package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"gitlab.com/Ash3941/project-management/services/issue-service/internal/models"
	"gitlab.com/Ash3941/project-management/services/issue-service/internal/store"
)

// Service contains issue-related business logic.
type Service struct{ store *store.Store }

// New returns a new issue service wired with the provided store.
func New(s *store.Store) *Service { return &Service{store: s} }

// CreateIssue creates a new issue tied to a project.
func (s *Service) CreateIssue(ctx context.Context, projectID, title, body string) (*models.Issue, error) {
	i := &models.Issue{ID: uuid.New().String(), ProjectID: projectID, Title: title, Body: body, Status: "open", Priority: "normal", CreatedAt: time.Now().UTC()}
	if err := s.store.CreateIssue(ctx, i); err != nil {
		return nil, err
	}
	return i, nil
}
