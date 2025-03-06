package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config represents the application configuration
type Config struct {
	JSONFolderPath string
	Port           string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Try to load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: Error loading .env file, using default values")
	}

	// Get JSON folder path from environment variable or use default
	jsonFolderPath := os.Getenv("JSON_FOLDER_PATH")
	if jsonFolderPath == "" {
		jsonFolderPath = "./endpoints"
		log.Printf("JSON_FOLDER_PATH not set, using default: %s", jsonFolderPath)
	}

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("PORT not set, using default: %s", port)
	}

	return &Config{
		JSONFolderPath: jsonFolderPath,
		Port:           port,
	}, nil
}
