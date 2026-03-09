package postgres

import (
	"context"

	"github.com/evgeney-fullstack/speed-reading-app-backend/internal/app/models"
	"github.com/jmoiron/sqlx"
)

// TextRepoStore defines CRUD operations for reading text management
type TextRepoStore interface {
	InsertText(ctx context.Context, text models.ReadingText) (int64, error)
	GetTextById(ctx context.Context, textID int64) (models.ReadingText, error)
	DeleteText(ctx context.Context, textID int64) error
	UpdateText(ctx context.Context, textID int64, input *models.UpdateReadingText) error
}

// Repository aggregates all store interfaces for database operations
type Repository struct {
	TextRepoStore
}

// NewRepository constructs a new Repository with all available stores
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		TextRepoStore: NewTextRepository(db),
	}
}
