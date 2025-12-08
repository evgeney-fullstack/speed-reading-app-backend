package postgres

import (
	"context"

	"github.com/evgeney-fullstack/speed-reading-app-backend/internal/app/models"
	"github.com/jmoiron/sqlx"
)

// TextRepoStore defines CRUD operations for reading text management
type TextRepoStore interface {
	InsertText(ctx context.Context, text models.ReadingText) (int64, error)
	GetAll()
	GetById()
	Delete()
	Update()
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
