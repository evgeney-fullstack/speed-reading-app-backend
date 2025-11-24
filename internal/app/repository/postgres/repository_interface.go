package postgres

import (
	"github.com/jmoiron/sqlx"
)

// SubscriptionStore defines CRUD operations for subscription management
type SubscriptionStore interface {
	Create()
	GetAll()
	GetById()
	Delete()
	Update()
}

// Repository aggregates all store interfaces for database operations
type Repository struct {
	SubscriptionStore
}

// NewRepository constructs a new Repository with all available stores
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		SubscriptionStore: NewSubscriptionRepository(db),
	}
}
