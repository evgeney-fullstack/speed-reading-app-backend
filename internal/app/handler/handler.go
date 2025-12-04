package handler

import (
	"github.com/evgeney-fullstack/speed-reading-app-backend/internal/app/service"
	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests and manages routing.
// Contains dependencies required for request handlers (future fields).
type Handler struct {
	services *service.Service
}

// NewHandler creates and returns a new Handler instance.
// Constructor function for initializing a handler with possible dependencies.
func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
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

		questions := speedReading.Group("/:text_id/questions")
		{
			questions.POST("/", h.createQuestion)               //Create a new question
			questions.GET("/", h.getAllQuestion)                //Retrieve all question
			questions.PUT("/:question_id", h.updateQuestion)    //Update an existing question
			questions.DELETE("/:question_id", h.deleteQuestion) //Delete a question
		}
	}

	return router

}
