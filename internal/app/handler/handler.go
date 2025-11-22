package handler

import (
	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests and manages routing.

// Contains dependencies required for request handlers (future fields).

type Handler struct {
}

// NewHandler creates and returns a new Handler instance.

// Constructor function for initializing a handler with possible dependencies.

func NewHandler() *Handler {

	return &Handler{}

}

// InitRoutes configures and returns the Gin router with defined endpoints.

// Adds middleware and registers handlers for all API paths.

func (h *Handler) InitRoutes() *gin.Engine {

	router := gin.New()

	return router

}
