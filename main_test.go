package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
		"path": "/users",
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
		"path": "/create-user",
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

	// Set environment variables - use a different port to avoid conflicts
	testPort := "8099"
	os.Setenv("JSON_FOLDER_PATH", tempDir)
	os.Setenv("PORT", testPort)
	defer func() {
		os.Unsetenv("JSON_FOLDER_PATH")
		os.Unsetenv("PORT")
	}()

	// Start the server
	srv, err := startTestServer()
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}

	// Make sure to shut down the server when the test is done
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Stop(ctx); err != nil {
			t.Errorf("Error shutting down server: %v", err)
		}
	}()

	// Wait for the server to start
	time.Sleep(1 * time.Second)

	// Test GET /users
	t.Run("GET /users", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%s/users", testPort))
		if err != nil {
			t.Fatalf("Failed to get users: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := ioutil.ReadAll(resp.Body)
			t.Errorf("Expected status OK, got %v with body: %s", resp.Status, string(bodyBytes))
			return
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

		resp, err := http.Post(fmt.Sprintf("http://localhost:%s/create-user", testPort), "application/json", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := ioutil.ReadAll(resp.Body)
			t.Errorf("Expected status Created, got %v with body: %s", resp.Status, string(bodyBytes))
			return
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

// startTestServer creates and starts a test server
func startTestServer() (*server.Server, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	mockResponses, err := mock.LoadResponses(cfg.JSONFolderPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load responses: %v", err)
	}

	srv, err := server.New(mockResponses, cfg.Port)
	if err != nil {
		return nil, fmt.Errorf("failed to create server: %v", err)
	}

	// Start server in a goroutine
	go func() {
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server error: %v\n", err)
		}
	}()

	return srv, nil
}
