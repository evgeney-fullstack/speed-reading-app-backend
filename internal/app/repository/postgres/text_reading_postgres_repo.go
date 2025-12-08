package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/evgeney-fullstack/speed-reading-app-backend/internal/app/models"
	"github.com/jmoiron/sqlx"
)

// TextServiceRepository implements reading text for PostgreSQL
type TextRepository struct {
	db *sqlx.DB
}

// NewTextServiceRepository creates a new instance of text repository
func NewTextRepository(db *sqlx.DB) *TextRepository {
	return &TextRepository{db: db}
}

// InsertText creates a new reading text record in the database
func (r *TextRepository) InsertText(ctx context.Context, text models.ReadingText) (int64, error) {
	// Check if context was cancelled before proceeding
	if err := ctx.Err(); err != nil {
		return 0, fmt.Errorf("context cancelled before database operation: %w", err)
	}

	// Set timestamps to current UTC time
	now := time.Now().UTC()

	// Only set CreatedAt if it hasn't been set (allows for data migration scenarios)
	if text.CreatedAt.IsZero() {
		text.CreatedAt = now
	}

	// Always update UpdatedAt on creation
	if text.UpdatedAt.IsZero() {
		text.UpdatedAt = now
	}

	// Build and execute SQL INSERT query
	query := fmt.Sprintf("INSERT INTO %s (content, word_count, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id", readingTextsTable)

	// Execute query and retrieve the auto-generated ID
	err := r.db.QueryRowContext(ctx, query,
		text.Content,
		text.WordCount,
		text.CreatedAt,
		text.UpdatedAt,
	).Scan(&text.ID)

	if err != nil {
		return 0, fmt.Errorf("failed to insert text into database: %w", err)
	}

	return text.ID, nil
}

// GetAll implements retrieval of all reading texts (to be implemented)
func (r *TextRepository) GetAll() {

}

// GetById implements retrieval of reading text by ID (to be implemented)
func (r *TextRepository) GetById() {

}

// Delete implements reading text deletion logic (to be implemented)
func (r *TextRepository) Delete() {

}

// Update implements reading text update logic (to be implemented)
func (r *TextRepository) Update() {

}
