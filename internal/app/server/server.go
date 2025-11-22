package server

import (
	"context"

	"net/http"

	"time"
)

// Server encapsulates an HTTP server and provides methods for managing its lifecycle.

type Server struct {
	httpServer *http.Server
}

// Run starts an HTTPS server with the specified host and port.

// Uses pre-configured TLS certificates (cert.pem and key.pem).

// Returns an error in case of failure (for example, problems downloading certificates or a busy port).

// Run starts an HTTPS server with the specified host and port.

// Uses pre-configured TLS certificates (cert.pem and key.pem).

// Returns an error in case of failure (for example, problems downloading certificates or a busy port).
func (s *Server) Run(host, port string, handler http.Handler) error {

	s.httpServer = &http.Server{

		Addr: host + ":" + port,

		Handler: handler,

		MaxHeaderBytes: 1 << 20,

		ReadTimeout: 10 * time.Second,

		WriteTimeout: 10 * time.Second,
	}

	// Launching a TLS-enabled server

	return s.httpServer.ListenAndServeTLS("cert.pem", "key.pem")

}

// Shutdown stops the server correctly, ending all active connections.

// Accepts a context for monitoring the execution time of the stop.

func (s *Server) Shutdown(ctx context.Context) error {

	return s.httpServer.Shutdown(ctx)

}
