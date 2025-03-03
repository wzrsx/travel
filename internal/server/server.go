package server

import (
	"fmt"
	"net/http"
	"travel/internal/config"
	"travel/internal/handlers"
	"travel/internal/repositories/pool_conections"
)

type Server struct {
	Host string
}

func NewServer() (*Server, error) {
	config, err := config.NewConfigIP()
	if(err != nil){
		return nil, err
	}
	return &Server{
		Host: config.AdressStr(),
	}, nil
}

func (s *Server) StartServe() error {
	pool, err := pool_conections.Create_pool()
	if err != nil {
		return err
	}

	app_handlers := handlers.NewAppHandlers(pool)
	app_handlers.SetHandlers()
	app_handlers.SetDirs()

	fmt.Printf("Server listening on port %s...\n", s.Host)
	if err := http.ListenAndServe(s.Host, nil); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}
