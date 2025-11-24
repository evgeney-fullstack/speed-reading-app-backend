package service

import "github.com/evgeney-fullstack/speed-reading-app-backend/internal/app/repository/postgres"

// SubscriptionService implements business logic for subscription operations
type SubscriptionService struct {
	repo postgres.SubscriptionStore
}

// NewSubscriptionService creates a new subscription service instance
func NewSubscriptionService(repo postgres.SubscriptionStore) *SubscriptionService {
	return &SubscriptionService{repo: repo}
}

// Create implements subscription creation business logic (to be implemented)
func (s *SubscriptionService) Create() {
}

// GetAll implements business logic for retrieving all subscriptions (to be implemented)
func (s *SubscriptionService) GetAll() {

}

// GetById implements business logic for retrieving subscription by ID (to be implemented)
func (s *SubscriptionService) GetById() {

}

// Delete implements subscription deletion business logic (to be implemented)
func (s *SubscriptionService) Delete() {

}

// Update implements subscription update business logic (to be implemented)
func (s *SubscriptionService) Update() {

}
