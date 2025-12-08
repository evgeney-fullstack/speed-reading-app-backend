package service

import (
	"context"
	"fmt"
	"strings"

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
	// Checking if the context has already been canceled.
	if err := ctx.Err(); err != nil {
		return 0, fmt.Errorf("context cancelled before operation: %w", err)
	}

	// Split by whitespace and filter out empty strings
	text.WordCount = len(strings.Fields(text.Content))

	return s.repo.InsertText(ctx, text)
}

// GetAll implements business logic for retrieving all reading texts (to be implemented)
func (s *TextService) GetAll() {

}

// GetById implements business logic for retrieving reading text by ID (to be implemented)
func (s *TextService) GetById() {

}

// Delete implements reading text deletion business logic (to be implemented)
func (s *TextService) Delete() {

}

// Update implements reading text update business logic (to be implemented)
func (s *TextService) Update() {

}
