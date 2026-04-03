package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/evgeney-fullstack/speed-reading-app-backend/internal/app/apperrors"
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

// updateReadingText updates the reading text by the specified identifier.
// Expects the text_id in the URL (a positive integer) and a JSON object with the fields to be updated.
// On success, it returns the 200 OK status with a success message.
//
// Possible errors:
// - 400: invalid ID format, ID <= 0, incorrect JSON, or missing update fields
// - 404: no text with this ID found in the database
// - 408: operation timeout (database operation exceeded 15 seconds)
// - 500: internal server error (unexpected issues during update)
func (h *Handler) updateReadingText(c *gin.Context) {
	// Extract and convert text_id parameter from URL path to integer
	// This parameter is required and must be a valid positive integer
	textID, err := strconv.Atoi(c.Param("text_id"))
	if err != nil {
		// Return 400 Bad Request if the parameter cannot be parsed as an integer
		// Example: passing "abc" or an empty string instead of a number
		h.errorHandler.BadRequest(c, "invalid_text_id_param", "Invalid text ID format: must be a valid integer")
		return
	}

	// Validate that the ID is positive (IDs <= 0 are not allowed)
	// Ensures we're working with a valid, existing resource identifier
	if textID <= 0 {
		h.errorHandler.BadRequest(c, "invalid_text_id_value", "Text ID must be a positive integer (greater than 0)")
		return
	}

	// Create context with timeout for database operation (15 seconds)
	// This prevents hanging operations and ensures timely responses
	ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
	defer cancel()

	var input models.UpdateReadingText
	if err := c.BindJSON(&input); err != nil {
		// Return 400 if JSON parsing fails
		// This covers malformed JSON or unexpected data types
		h.errorHandler.BadRequest(c, "invalid_json", "Invalid request body: malformed or unexpected JSON structure")
		return
	}

	// Check if at least one field to update is provided
	// Both Content and Questions being nil means no update data was sent
	if input.Content == nil && input.Questions == nil {
		// Return 400 with specific code and message
		// Clarifies that the request is missing required update data
		h.errorHandler.BadRequest(
			c,
			"missing_update_fields",
			"At least one of the fields 'content' or 'questions' must be provided for update. Empty updates are not allowed",
		)
		return
	}

	// Attempt to update the reading text using the service layer
	// Passes the context (with timeout), text ID, and update data
	err = h.services.TextServiceStore.UpdateReadingText(ctx, int64(textID), input)
	if err != nil {
		// Check specifically for context timeout errors
		// Indicates the database operation took longer than 15 seconds
		if errors.Is(err, context.DeadlineExceeded) {
			h.errorHandler.RequestTimeout(c, "request_timeout", "Operation timed out: update took longer than 15 seconds")
			return
		}

		// Handle case where the requested text ID doesn't exist in the database
		// Returns 404 Not Found to indicate resource absence
		if errors.Is(err, apperrors.ErrTextNotFound) {
			h.errorHandler.NotFound(c, "Text not found: no record exists with the specified ID")
			return
		}

		// For any other unexpected errors, return 500 Internal Server Error
		// Logs detailed error internally but returns generic message to client
		// Prevents exposing sensitive internal details to API consumers
		h.errorHandler.InternalError(c, err, "Failed to update reading text due to an internal server error")
		return
	}

	// If all steps succeed, return 200 OK with success message
	// Indicates the update operation completed successfully
	h.errorHandler.Success(c, "ok")
}

// DeleteReadingText removes a reading text by its ID
// DELETE /reading_text/:text_id
func (h *Handler) deleteReadingText(c *gin.Context) {
	// Extract and convert text_id parameter from URL path to integer
	textID, err := strconv.Atoi(c.Param("text_id"))
	if err != nil {
		// Return 400 Bad Request if the parameter is not a valid integer
		h.errorHandler.BadRequest(c, "invalid_text_id_param", "Invalid text ID format")
		return
	}

	// Validate for positive ID value
	if textID <= 0 {
		h.errorHandler.BadRequest(c, "invalid_text_id_value", "Text ID must be positive")
		return
	}

	// Create context with timeout for database operation
	ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
	defer cancel()

	// Delete the reading text from the service layer
	err = h.services.TextServiceStore.DeleteReadingText(ctx, int64(textID))
	if err != nil {
		// Check for context timeout
		if errors.Is(err, context.DeadlineExceeded) {
			h.errorHandler.RequestTimeout(c, "request_timeout", "Operation timed out")
			return
		}

		// Handle "not found" error from service layer
		if errors.Is(err, apperrors.ErrTextNotFound) {
			h.errorHandler.NotFound(c, "Text not found")
			return
		}

		// Log internal error details but return generic message to client
		h.errorHandler.InternalError(c, err, "Failed to delete reading text")
		return
	}

	// Return HTTP 204 OK with success confirmation
	c.Status(http.StatusNoContent)
}
