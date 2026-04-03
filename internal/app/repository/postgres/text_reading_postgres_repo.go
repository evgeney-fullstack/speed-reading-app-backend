package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/evgeney-fullstack/speed-reading-app-backend/internal/app/apperrors"
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
	query := fmt.Sprintf("INSERT INTO %s (content, word_count, questions, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id", readingTextsTable)

	// Execute query and retrieve the auto-generated ID
	err := r.db.QueryRowContext(ctx, query,
		text.Content,
		text.WordCount,
		text.Questions,
		text.CreatedAt,
		text.UpdatedAt,
	).Scan(&text.ID)

	if err != nil {
		return 0, fmt.Errorf("failed to insert text into database: %w", err)
	}

	return text.ID, nil
}

// GetTextById retrieves a reading text by ID from the database
// Returns sql.ErrNoRows if no text found with the given ID
func (r *TextRepository) GetTextById(ctx context.Context, textID int64) (models.ReadingText, error) {
	var text models.ReadingText

	// Check if context was cancelled before proceeding with database operation
	if err := ctx.Err(); err != nil {
		return text, fmt.Errorf("context cancelled before database query: %w", err)
	}

	// Using parameterized query to prevent SQL injection
	// readingTextsTable is a constant defined elsewhere
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", readingTextsTable)

	// Execute the query with context for proper timeout/cancellation handling
	err := r.db.GetContext(ctx, &text, query, textID)
	if err != nil {
		// Return error as-is to allow service layer to handle specific cases
		// (e.g., sql.ErrNoRows for "not found" scenario)
		return text, err
	}

	return text, nil
}

// DeleteText removes a reading text by ID from the database
// Returns an error if the text is not found or if operation fails
func (r *TextRepository) DeleteText(ctx context.Context, textID int64) error {
	// Check if context was cancelled before proceeding
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("context cancelled before delete operation: %w", err)
	}

	// Using parameterized query to prevent SQL injection
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", readingTextsTable)

	// Execute the delete operation
	result, err := r.db.ExecContext(ctx, query, textID)
	if err != nil {
		return fmt.Errorf("database error during delete: %w", err)
	}

	// Check if any row was actually deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%w: text with ID %d not found", apperrors.ErrTextNotFound, textID)
	}

	return nil
}

// Update implements reading text update logic (to be implemented)
func (r *TextRepository) UpdateText(ctx context.Context, textID int64, input *models.UpdateReadingText) error {

	// Initialize slices for building dynamic SET clause and arguments
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1 // Positional parameter counter

	// Handle price update if provided
	if input.Content != nil {
		setValues = append(setValues, fmt.Sprintf("content=$%d", argId))
		args = append(args, *input.Content)
		argId++

		setValues = append(setValues, fmt.Sprintf("word_count=$%d", argId))
		args = append(args, input.WordCount)
		argId++
	}

	// Handle start date update if provided
	if input.Questions != nil {
		setValues = append(setValues, fmt.Sprintf("questions=$%d", argId))
		args = append(args, *input.Questions)
		argId++
	}

	updatedAt := time.Now().UTC()
	setValues = append(setValues, fmt.Sprintf("updated_at=$%d", argId))
	args = append(args, updatedAt)
	argId++

	// Join SET clauses with commas
	setQuery := strings.Join(setValues, ", ")

	// Build final SQL query with WHERE clause
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d", readingTextsTable, setQuery, argId)

	// Add subscription ID as the last parameter
	args = append(args, textID)

	// Execute the query
	_, err := r.db.Exec(query, args...)
	return err
}
