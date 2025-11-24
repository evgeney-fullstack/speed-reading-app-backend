package service

import "github.com/evgeney-fullstack/speed-reading-app-backend/internal/app/repository/postgres"

// SubscriptionStore defines business logic operations for subscriptions
type SubscriptionStore interface {
	Create()
	GetAll()
	GetById()
	Delete()
	Update()
}

// Service layer aggregates all business logic services
type Service struct {
	SubscriptionStore
}

// NewService constructs new Service layer with business logic
func NewService(repos *postgres.Repository) *Service {
	return &Service{
		SubscriptionStore: NewSubscriptionService(repos.SubscriptionStore),
	}
}
