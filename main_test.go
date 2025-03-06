package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sachin-duhan/gomock/pkg/config"
	"github.com/sachin-duhan/gomock/pkg/mock"
	"github.com/sachin-duhan/gomock/pkg/server"
)

func TestIntegration(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := ioutil.TempDir("", "gomock-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test JSON files
	usersJSON := `{
		"method": "GET",
		"responses": [
			{
				"status": 200,
				"body": {"users": [{"id": 1, "name": "Test User"}]},
				"description": "Success response"
			}
		]
	}`

	createUserJSON := `{
		"method": "POST",
		"responses": [
			{
				"status": 201,
				"body": {"message": "Success"},
				"input_body": {"name": "Test User", "email": "test@example.com"},
				"description": "Success case"
			}
		]
	}`

	// Write test files
	if err := ioutil.WriteFile(filepath.Join(tempDir, "users.json"), []byte(usersJSON), 0644); err != nil {
		t.Fatalf("Failed to write users.json: %v", err)
	}
	if err := ioutil.WriteFile(filepath.Join(tempDir, "create-user.json"), []byte(createUserJSON), 0644); err != nil {
		t.Fatalf("Failed to write create-user.json: %v", err)
	}

	// Set environment variables
	os.Setenv("JSON_FOLDER_PATH", tempDir)
	os.Setenv("PORT", "8081")
	defer func() {
		os.Unsetenv("JSON_FOLDER_PATH")
		os.Unsetenv("PORT")
	}()

	// Start the server in a goroutine
	go func() {
		if err := run(); err != nil {
			t.Errorf("Server error: %v", err)
		}
	}()

	// Wait for the server to start
	time.Sleep(time.Second)

	// Test GET /users
	t.Run("GET /users", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8081/users")
		if err != nil {
			t.Fatalf("Failed to get users: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK, got %v", resp.Status)
		}

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		users, ok := result["users"].([]interface{})
		if !ok {
			t.Fatal("Expected users array in response")
		}

		if len(users) != 1 {
			t.Errorf("Expected 1 user, got %d", len(users))
		}
	})

	// Test POST /create-user
	t.Run("POST /create-user", func(t *testing.T) {
		body := map[string]string{
			"name":  "Test User",
			"email": "test@example.com",
		}
		jsonBody, _ := json.Marshal(body)

		resp, err := http.Post("http://localhost:8081/create-user", "application/json", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected status Created, got %v", resp.Status)
		}

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if msg, ok := result["message"].(string); !ok || msg != "Success" {
			t.Errorf("Expected success message, got %v", result)
		}
	})
}

func run() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	mockResponses, err := mock.LoadResponses(cfg.JSONFolderPath)
	if err != nil {
		return err
	}

	srv, err := server.New(mockResponses, cfg.Port)
	if err != nil {
		return err
	}

	return srv.Start()
}
