package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Save original environment variables
	originalPort := os.Getenv("PORT")
	originalJSONFolderPath := os.Getenv("JSON_FOLDER_PATH")

	// Restore environment variables after the test
	defer func() {
		os.Setenv("PORT", originalPort)
		os.Setenv("JSON_FOLDER_PATH", originalJSONFolderPath)
	}()

	// Test case 1: Default values when environment variables are not set
	os.Unsetenv("PORT")
	os.Unsetenv("JSON_FOLDER_PATH")

	config := LoadConfig()

	if config.Port != "8080" {
		t.Errorf("Expected default port 8080, got %s", config.Port)
	}

	if config.JSONFolderPath != "./endpoints" {
		t.Errorf("Expected default JSON folder path ./endpoints, got %s", config.JSONFolderPath)
	}

	// Test case 2: Custom values from environment variables
	os.Setenv("PORT", "9090")
	os.Setenv("JSON_FOLDER_PATH", "./custom-endpoints")

	config = LoadConfig()

	if config.Port != "9090" {
		t.Errorf("Expected port 9090, got %s", config.Port)
	}

	if config.JSONFolderPath != "./custom-endpoints" {
		t.Errorf("Expected JSON folder path ./custom-endpoints, got %s", config.JSONFolderPath)
	}
}
