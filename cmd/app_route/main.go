package main

import (
	"log"
	"travel/internal/server"
)

func main() {
	srv := server.NewServer()
	err := srv.StartServe()
	if err != nil{
		log.Fatalf("Error start server: %s", err.Error())
		return
	}
}
