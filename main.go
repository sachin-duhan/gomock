package main

import (
	"log"

	"github.com/sachin-duhan/gomock/pkg/config"
	"github.com/sachin-duhan/gomock/pkg/mock"
	"github.com/sachin-duhan/gomock/pkg/server"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Load mock responses
	mockResponses, err := mock.LoadResponses(cfg.JSONFolderPath)
	if err != nil {
		log.Fatalf("Failed to load mock responses: %v", err)
	}

	// Create and start the server
	srv, err := server.New(mockResponses, cfg.Port)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	if err := srv.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
