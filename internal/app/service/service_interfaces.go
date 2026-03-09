package service

import (
	"context"

	"github.com/evgeney-fullstack/speed-reading-app-backend/internal/app/models"
	"github.com/evgeney-fullstack/speed-reading-app-backend/internal/app/repository/postgres"
)

// TextServiceStore defines business logic operations for reading text
type TextServiceStore interface {
	CreateReadingText(ctx context.Context, text models.ReadingText) (int64, error)
	GetReadingText(ctx context.Context, textID int64) (models.ReadingText, error)
	DeleteReadingText(ctx context.Context, textID int64) error
	UpdateReadingText(ctx context.Context, textID int64, input models.UpdateReadingText) error
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
