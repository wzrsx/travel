package main

import (
	"log"
	"travel/internal/server"
)

func main() {
	srv, err := server.NewServer()
	if err != nil{
		log.Fatalf("Error get env temporaries: %s", err.Error())
	}
	err = srv.StartServe()
	if err != nil{
		log.Fatalf("Error start server: %s", err.Error())
		return
	}
}
