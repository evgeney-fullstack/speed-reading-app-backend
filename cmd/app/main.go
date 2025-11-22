package main

import (
	"log"

	"github.com/evgeney-fullstack/speed-reading-app-backend/internal/app/handler"
	"github.com/evgeney-fullstack/speed-reading-app-backend/internal/app/server"
)

func main() {

	handlers := handler.NewHandler()

	srv := new(server.Server)

	//Running an HTTP server from TLS to localhost:8080

	if err := srv.Run("localhost", "8080", handlers.InitRoutes()); err != nil {

		log.Fatalf("error occurred while running http server: %s", err.Error())

	}

}
