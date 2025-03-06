package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Test with environment variables set
	t.Run("With environment variables", func(t *testing.T) {
		// Set environment variables
		os.Setenv("JSON_FOLDER_PATH", "/test/path")
		os.Setenv("PORT", "3000")
		defer func() {
			os.Unsetenv("JSON_FOLDER_PATH")
			os.Unsetenv("PORT")
		}()

		// Load config
		cfg, err := LoadConfig()
		if err != nil {
			t.Fatalf("LoadConfig failed: %v", err)
		}

		// Check values
		if cfg.JSONFolderPath != "/test/path" {
			t.Errorf("Expected JSONFolderPath to be /test/path, got %s", cfg.JSONFolderPath)
		}
		if cfg.Port != "3000" {
			t.Errorf("Expected Port to be 3000, got %s", cfg.Port)
		}
	})

	// Test with default values
	t.Run("With default values", func(t *testing.T) {
		// Load config without environment variables
		cfg, err := LoadConfig()
		if err != nil {
			t.Fatalf("LoadConfig failed: %v", err)
		}

		// Check default values
		if cfg.JSONFolderPath != "./endpoints" {
			t.Errorf("Expected default JSONFolderPath to be ./endpoints, got %s", cfg.JSONFolderPath)
		}
		if cfg.Port != "8080" {
			t.Errorf("Expected default Port to be 8080, got %s", cfg.Port)
		}
	})
}
