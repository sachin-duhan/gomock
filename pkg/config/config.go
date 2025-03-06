package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config represents the application configuration
type Config struct {
	Port           string
	JSONFolderPath string
}

// LoadConfig loads the application configuration from environment variables
func LoadConfig() *Config {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file, using default values")
	}

	// Get folder path from environment variable
	folderPath := os.Getenv("JSON_FOLDER_PATH")
	if folderPath == "" {
		folderPath = "./endpoints" // Default value
		log.Printf("JSON_FOLDER_PATH not set, using default: %s", folderPath)
	}

	// Get server port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default value
		log.Printf("PORT not set, using default: %s", port)
	}

	return &Config{
		Port:           port,
		JSONFolderPath: folderPath,
	}
}
