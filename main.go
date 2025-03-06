package main

import (
	"log"

	"github.com/sachin-duhan/gomock/pkg/config"
	"github.com/sachin-duhan/gomock/pkg/mock"
	"github.com/sachin-duhan/gomock/pkg/server"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Load mock responses from JSON files
	mockResponses, err := mock.LoadResponses(cfg.JSONFolderPath)
	if err != nil {
		log.Fatalf("Error loading mock responses: %v", err)
	}

	// Log loaded endpoints
	log.Printf("Loaded %d mock endpoints", len(mockResponses))
	for endpoint, mock := range mockResponses {
		log.Printf("Endpoint: %s, Method: %s, Status: %d",
			endpoint, mock.Method, mock.Response.Status)
	}

	// Create and start the server
	srv := server.New(mockResponses, cfg.Port)
	log.Fatal(srv.Start())
}
