package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/evgeney-fullstack/speed-reading-app-backend/internal/app/apperrors"
	"github.com/evgeney-fullstack/speed-reading-app-backend/internal/app/models"
	"github.com/evgeney-fullstack/speed-reading-app-backend/internal/app/repository/postgres"
)

// TextService implements business logic for text operations
type TextService struct {
	repo postgres.TextRepoStore
}

// NewTextService creates a new instance of text service
func NewTextService(repo postgres.TextRepoStore) *TextService {
	return &TextService{repo: repo}
}

// CreateReadingText implements reading text creation business logic (to be implemented)
func (s *TextService) CreateReadingText(ctx context.Context, text models.ReadingText) (int64, error) {
	// Check if context is still valid before proceeding
	if err := ctx.Err(); err != nil {
		return 0, fmt.Errorf("context error before repository call: %w", err)
	}

	// Split by whitespace and filter out empty strings
	text.WordCount = len(strings.Fields(text.Content))

	return s.repo.InsertText(ctx, text)
}

// GetReadingText retrieves a reading text by ID from the repository
// Performs business logic validation and error handling
func (s *TextService) GetReadingText(ctx context.Context, textID int64) (models.ReadingText, error) {

	// Check if context is still valid before proceeding
	if err := ctx.Err(); err != nil {
		return models.ReadingText{}, fmt.Errorf("context error before repository call: %w", err)
	}

	// Retrieve reading text by ID from the repository layer (database)
	text, err := s.repo.GetTextById(ctx, textID)
	if err != nil {
		// Wrap repository error with service layer context
		return models.ReadingText{}, fmt.Errorf("failed to retrieve reading text with ID %d: %w", textID, err)
	}

	return text, nil
}

// DeleteReadingText handles business logic for deleting a reading text by ID
// Validates input and delegates to repository layer
func (s *TextService) DeleteReadingText(ctx context.Context, textID int64) error {

	// Check if context is still valid before proceeding
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("context error before delete operation: %w", err)
	}

	// Delegate deletion to repository layer
	err := s.repo.DeleteText(ctx, textID)
	if err != nil {

		//Service  "not found" error from repository layer
		if errors.Is(err, apperrors.ErrTextNotFound) {
			return apperrors.ErrTextNotFound
		}

		// Wrap repository error with service layer context
		return fmt.Errorf("failed to delete reading text with ID %d: %w", textID, err)
	}

	return nil
}

// Update implements reading text update business logic (to be implemented)
func (s *TextService) Update() {

}
