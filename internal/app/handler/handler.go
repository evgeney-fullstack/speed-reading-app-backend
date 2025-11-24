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

	return router

}
