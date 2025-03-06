package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/sachin-duhan/gomock/pkg/config"
	"github.com/sachin-duhan/gomock/pkg/mock"
)

func TestIntegration(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := ioutil.TempDir("", "gomock-integration-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test JSON files
	testFiles := []struct {
		filename string
		content  string
	}{
		{
			filename: "users.json",
			content: `{
				"method": "GET",
				"response": {
					"status": 200,
					"body": {"users": [{"id": 1, "name": "Test User"}]}
				}
			}`,
		},
		{
			filename: "create-user.json",
			content: `{
				"method": "POST",
				"response": {
					"status": 201,
					"body": {"message": "User created", "id": 2}
				}
			}`,
		},
	}

	for _, tf := range testFiles {
		filePath := filepath.Join(tempDir, tf.filename)
		err := ioutil.WriteFile(filePath, []byte(tf.content), 0644)
		if err != nil {
			t.Fatalf("Failed to write test file %s: %v", tf.filename, err)
		}
	}

	// Save original environment variables
	originalPort := os.Getenv("PORT")
	originalJSONFolderPath := os.Getenv("JSON_FOLDER_PATH")

	// Set environment variables for the test
	os.Setenv("PORT", "8081")
	os.Setenv("JSON_FOLDER_PATH", tempDir)

	// Restore environment variables after the test
	defer func() {
		os.Setenv("PORT", originalPort)
		os.Setenv("JSON_FOLDER_PATH", originalJSONFolderPath)
	}()

	// Test loading configuration
	cfg := config.LoadConfig()

	if cfg.Port != "8081" {
		t.Errorf("Expected port 8081, got %s", cfg.Port)
	}

	if cfg.JSONFolderPath != tempDir {
		t.Errorf("Expected JSON folder path %s, got %s", tempDir, cfg.JSONFolderPath)
	}

	// Test loading mock responses
	mockResponses, err := mock.LoadResponses(cfg.JSONFolderPath)
	if err != nil {
		t.Fatalf("Failed to load mock responses: %v", err)
	}

	if len(mockResponses) != len(testFiles) {
		t.Errorf("Expected %d mock responses, got %d", len(testFiles), len(mockResponses))
	}

	// Check if all endpoints are loaded
	expectedEndpoints := []string{"/users", "/create-user"}
	for _, endpoint := range expectedEndpoints {
		if _, exists := mockResponses[endpoint]; !exists {
			t.Errorf("Expected endpoint %s not found in loaded responses", endpoint)
		}
	}
}
