package server

import (
	"fmt"
	"log"
	"net/http"
	"travel/internal/config"
	"travel/internal/handlers"
	"travel/internal/http_handlers"
)

type Server struct {
	Host string
}

func NewServer() *Server {
	return &Server{
		Host: config.NewConfigIP().AdressStr(),
	}
}

func (s *Server) StartServe() error {
	http.Handle("/", handlers.LoggingMiddleware(http.HandlerFunc(http_handlers.OpenFirstPage)))

	http_handlers.SetDirs()

	fmt.Printf("Server listening on port %s...\n", s.Host)
	if err := http.ListenAndServe(s.Host, nil); err != nil {
		log.Printf("Server failed to start: %v", err)
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}
