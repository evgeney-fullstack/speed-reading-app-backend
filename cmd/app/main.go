package main

import (
	"os"

	"github.com/evgeney-fullstack/speed-reading-app-backend/internal/app/handler"
	"github.com/evgeney-fullstack/speed-reading-app-backend/internal/app/server"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	// Configuring the logs format in JSON for better structuring and compatibility
	// with monitoring systems (Kibana, Elasticsearch, etc.)
	logrus.SetFormatter(new(logrus.JSONFormatter))

	// Loading environment variables from the config.env file
	if err := godotenv.Load("config.env"); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	// Initialization of HTTP request handlers
	handlers := handler.NewHandler()

	// Creating a server instance
	srv := new(server.Server)

	// Launching an HTTPS server with configuration from environment variables
	// Using HOST and HOST_PORT from config.env
	if err := srv.Run(os.Getenv("HOST"), os.Getenv("HOST_PORT"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error occurred while running http server: %s", err.Error())
	}

}
