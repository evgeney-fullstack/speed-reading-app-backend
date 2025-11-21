package main

import (
	"log"

	"github.com/evgeney-fullstack/speed-reading-app-backend/internal/app/server"
)

func main() {
	srv := new(server.Server)

	//Running an HTTP server from TLS to localhost:8080

	if err := srv.Run("localhost", "8080"); err != nil {

		log.Fatalf("error occurred while running http server: %s", err.Error())

	}

}
