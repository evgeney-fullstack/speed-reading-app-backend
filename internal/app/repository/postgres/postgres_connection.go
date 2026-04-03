package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	readingTextsTable = "reading_texts" // Database table name for reading_texts
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

// NewPostgresDB creates a new PostgreSQL database connection with context support
// Returns sqlx.DB instance or error if connection fails
func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)

	var db *sqlx.DB
	var err error

	// Retry logic for connection
	for i := 0; i < 3; i++ {
		db, err = sqlx.Open("postgres", connectionString)
		if err != nil {
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}

		if err = db.PingContext(ctx); err != nil {
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}

		break
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect after retries: %w", err)
	}

	return db, nil
}

// RunMigrations applies database migrations from the migrations folder.
// Returns error if migration fails (ignores the case when there are no changes).
func RunMigrations(db *sqlx.DB) error {
	// Initialize PostgreSQL driver instance
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return err
	}

	// Create migrator instance with file-based migrations
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations", // migrations are loaded from filesystem
		"postgres", driver)
	if err != nil {
		return err
	}

	// Apply all pending migrations (Up)
	err = m.Up()
	// Ignore "no change" error (when all migrations are already applied)
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	// Log success and return
	log.Println("Migrations applied successfully")
	return nil
}
