package handler

import (
	"github.com/evgeney-fullstack/speed-reading-app-backend/internal/app/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Handler handles HTTP requests and manages routing.
// Contains dependencies required for request handlers (future fields).
type Handler struct {
	services     *service.Service
	errorHandler *ErrorHandler
}

// NewHandler creates and returns a new Handler instance.
// Constructor function for initializing a handler with possible dependencies.
func NewHandler(services *service.Service, logger *logrus.Logger) *Handler {
	return &Handler{
		services:     services,
		errorHandler: NewErrorHandler(logger),
	}

}

// InitRoutes configures and returns the Gin router with defined endpoints.
// Adds middleware and registers handlers for all API paths.
func (h *Handler) InitRoutes() *gin.Engine {

	router := gin.New()

	// Create a route group for speedReading-related endpoints
	speedReading := router.Group("/reading_text")
	{
		speedReading.POST("/", h.createReadingText)           //Create a new reading text
		speedReading.GET("/:text_id", h.getReadingTextById)   //Get a specific reading text by ID
		speedReading.PUT("/:text_id", h.updateReadingText)    //Update an existing reading text
		speedReading.DELETE("/:text_id", h.deleteReadingText) //Delete a reading text

	}

	return router

}
