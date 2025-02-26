package main

import (
	"travel/internal/server"
)

func main() {
	srv := server.NewServer()
	srv.StartServe()
}
