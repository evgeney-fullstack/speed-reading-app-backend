package service

import (
	"context"

	"github.com/evgeney-fullstack/speed-reading-app-backend/internal/app/models"
	"github.com/evgeney-fullstack/speed-reading-app-backend/internal/app/repository/postgres"
)

// TextServiceStore defines business logic operations for reading text
type TextServiceStore interface {
	CreateReadingText(ctx context.Context, text models.ReadingText) (int64, error)
	GetAll()
	GetById()
	Delete()
	Update()
}

// Service layer aggregates all business logic services
type Service struct {
	TextServiceStore
}

// NewService constructs new Service layer with business logic
func NewService(repos *postgres.Repository) *Service {
	return &Service{
		TextServiceStore: NewTextService(repos.TextRepoStore),
	}
}
