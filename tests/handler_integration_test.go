package tests

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/evgeney-fullstack/speed-reading-app-backend/internal/app/handler"
	"github.com/evgeney-fullstack/speed-reading-app-backend/internal/app/repository/postgres"
	"github.com/evgeney-fullstack/speed-reading-app-backend/internal/app/service"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TestMain installs the test environment
func TestMain(m *testing.M) {
	// Setting Gin mode to test mode
	gin.SetMode(gin.TestMode)

	// Running the tests
	m.Run()
}

// setupTestContainer configures the PostgreSQL container for testing
func setupTestContainer(ctx context.Context) (postgres.Config, func(), error) {
	// Launching a PostgreSQL container using GenericContainer
	postgresReq := testcontainers.ContainerRequest{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").WithStartupTimeout(30 * time.Second),
	}

	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: postgresReq,
		Started:          true,
	})
	if err != nil {
		return postgres.Config{}, nil, fmt.Errorf("failed to start PostgreSQL container: %w", err)
	}

	// Getting the PostgreSQL connection parameters
	host, err := postgresContainer.Host(ctx)
	if err != nil {
		return postgres.Config{}, nil, err
	}
	port, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		return postgres.Config{}, nil, err
	}

	// Configuration for connecting to the PostgreSQL test container
	postgresCfg := postgres.Config{
		Host:     host,
		Port:     port.Port(),
		Username: "testuser",
		Password: "testpass",
		DBName:   "testdb",
		SSLMode:  "disable",
	}

	// Cleaning function to stop the container
	cleanup := func() {
		if terminateErr := postgresContainer.Terminate(ctx); terminateErr != nil {
			log.Printf("failed to terminate PostgreSQL container: %v", terminateErr)
		}
	}

	return postgresCfg, cleanup, nil
}

// setupTestServer creates and configures a test server
func setupTestServer(postgresCfg postgres.Config) (*gin.Engine, error) {

	// Initializing the logger
	logger := logrus.New()

	// Setting the logging level
	logger.SetLevel(logrus.InfoLevel)

	// Configuring the logs format in JSON for better structuring and compatibility
	// with monitoring systems (Kibana, Elasticsearch, etc.)
	logger.SetFormatter(new(logrus.JSONFormatter))

	// Initializing PostgreSQL
	db, err := postgres.NewPostgresDB(postgresCfg)
	if err != nil {
		return nil, err
	}

	// Add table creation
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS reading_texts (
    id BIGSERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    word_count INTEGER NOT NULL,
    questions JSONB NOT NULL DEFAULT '[]',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
        )
    `)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	// Initializing repositories
	repos := postgres.NewRepository(db)

	// Initialization of services
	services := service.NewService(repos)

	// Initializing handlers
	handler := handler.NewHandler(services, logger)

	// Setting up routes
	router := handler.InitRoutes()

	return router, nil
}
