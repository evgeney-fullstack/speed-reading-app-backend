package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"
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

	// Create context with timeout for database operation
	ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
	defer cancel()

	// Call service layer with context
	textID, err := h.services.TextServiceStore.CreateReadingText(ctx, readingText)
	if err != nil {
		// Check for specific database errors
		if errors.Is(err, context.DeadlineExceeded) {
			h.errorHandler.RequestTimeout(c, "request_timeout", "Operation timed out")
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

// GetReadingTextById retrieves a reading text by its ID
// GET /reading_text/:text_id
func (h *Handler) getReadingTextById(c *gin.Context) {
	// Extract and convert text_id parameter from URL path to integer
	// The parameter is expected to be in the format: /reading_text/{text_id}
	textID, err := strconv.Atoi(c.Param("text_id"))
	if err != nil {
		// Return 400 Bad Request if the parameter is not a valid integer
		h.errorHandler.BadRequest(c, "invalid text_id param", "Invalid request param format")
		return
	}

	//Validation for negative and null ID
	if textID <= 0 {
		h.errorHandler.BadRequest(c, "invalid_text_id_value", "Text ID must be positive")
		return
	}

	// Create context with timeout for database operation
	ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
	defer cancel()

	// Retrieve the reading text from the service layer using the extracted ID
	readingText, err := h.services.TextServiceStore.GetReadingText(ctx, int64(textID))
	if err != nil {

		// Check for specific database errors
		if errors.Is(err, context.DeadlineExceeded) {
			h.errorHandler.RequestTimeout(c, "request_timeout", "Operation timed out")
			return
		}

		// Log internal error details but return generic message to client
		h.errorHandler.InternalError(c, err, "Failed to get reading text")
		return
	}

	// Return 200 OK with the reading text data in JSON format
	c.JSON(http.StatusOK, readingText)
}

func (h *Handler) updateReadingText(c *gin.Context) {

}

func (h *Handler) deleteReadingText(c *gin.Context) {

}
