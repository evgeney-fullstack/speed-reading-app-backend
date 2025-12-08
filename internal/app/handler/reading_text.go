package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/evgeney-fullstack/speed-reading-app-backend/internal/app/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) createReadingText(c *gin.Context) {
	// Bind JSON request body to ReadingText model
	// Validates required fields based on 'binding' tags in the model
	var readingText models.ReadingText
	if err := c.BindJSON(&readingText); err != nil {
		// Return 400 Bad Request if JSON is malformed or validation fails
		h.errorHandler.BadRequest(c, "invalid_request_body", "Invalid request body format")

		return
	}

	// Validate that content is not empty
	if readingText.Content == "" {
		h.errorHandler.BadRequest(c, "empty_content", "Text content cannot be empty")
		return
	}

	// Create context with timeout for database operation
	ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
	defer cancel()

	// Call service layer with context
	textID, err := h.services.TextServiceStore.CreateReadingText(ctx, readingText)
	if err != nil {
		// Check for specific database errors
		if errors.Is(err, context.DeadlineExceeded) {
			h.errorHandler.BadRequest(c, "request_timeout", "Operation timed out")
			return
		}

		// Log internal error details but return generic message to client
		h.errorHandler.InternalError(c, err, "Failed to create reading text")
		return
	}

	// Return 201 Created with the ID of the newly created text
	c.JSON(http.StatusCreated, map[string]interface{}{
		"id": textID,
	})

}

func (h *Handler) getReadingTextById(c *gin.Context) {

}

func (h *Handler) updateReadingText(c *gin.Context) {

}

func (h *Handler) deleteReadingText(c *gin.Context) {

}
