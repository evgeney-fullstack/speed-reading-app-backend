package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	subscriptionTable = "subscription" // Database table name for subscriptions
)

// Config holds PostgreSQL connection configuration parameters
type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

// NewPostgresDB creates a new PostgreSQL database connection
// Returns sqlx.DB instance or error if connection fails
func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	// Format connection string from configuration parameters
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := sqlx.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Verify connection with ping
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	return db, nil
}
