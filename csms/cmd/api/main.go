package main

import (
	"log"

	"github.com/malikkhoiri/csms/internal/config"
	"github.com/malikkhoiri/csms/internal/server"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	srv, err := server.NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	srv.SetupRoutes()

	log.Fatal(srv.Start())
}
