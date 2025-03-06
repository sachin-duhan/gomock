package mock

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadResponses(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := ioutil.TempDir("", "gomock-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test JSON files
	testCases := []struct {
		filename string
		content  string
		expected Response
	}{
		{
			filename: "users.json",
			content: `{
				"method": "GET",
				"responses": [
					{
						"status": 200,
						"body": {"users": [{"id": 1, "name": "Test User"}]},
						"description": "Success response"
					},
					{
						"status": 401,
						"body": {"error": "Unauthorized"},
						"description": "Error response"
					}
				]
			}`,
			expected: Response{
				Method: "GET",
				Responses: []ResponseConfig{
					{
						Status:      200,
						Body:        map[string]interface{}{"users": []interface{}{map[string]interface{}{"id": float64(1), "name": "Test User"}}},
						Description: "Success response",
					},
					{
						Status:      401,
						Body:        map[string]interface{}{"error": "Unauthorized"},
						Description: "Error response",
					},
				},
			},
		},
		{
			filename: "create-user.json",
			content: `{
				"method": "POST",
				"responses": [
					{
						"status": 201,
						"body": {"message": "User created", "id": 2},
						"input_body": {"name": "Test User", "email": "test@example.com"},
						"description": "Success response"
					},
					{
						"status": 400,
						"body": {"error": "Invalid input"},
						"input_body": {"name": "Test User"},
						"description": "Error response"
					}
				]
			}`,
			expected: Response{
				Method: "POST",
				Responses: []ResponseConfig{
					{
						Status:      201,
						Body:        map[string]interface{}{"message": "User created", "id": float64(2)},
						InputBody:   map[string]interface{}{"name": "Test User", "email": "test@example.com"},
						Description: "Success response",
					},
					{
						Status:      400,
						Body:        map[string]interface{}{"error": "Invalid input"},
						InputBody:   map[string]interface{}{"name": "Test User"},
						Description: "Error response",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		filePath := filepath.Join(tempDir, tc.filename)
		err := ioutil.WriteFile(filePath, []byte(tc.content), 0644)
		if err != nil {
			t.Fatalf("Failed to write test file %s: %v", tc.filename, err)
		}
	}

	// Test loading mock responses
	mockResponses, err := LoadResponses(tempDir)
	if err != nil {
		t.Fatalf("LoadResponses failed: %v", err)
	}

	// Verify the loaded responses
	if len(mockResponses) != len(testCases) {
		t.Errorf("Expected %d mock responses, got %d", len(testCases), len(mockResponses))
	}

	for _, tc := range testCases {
		endpoint := "/" + filepath.Base(tc.filename[:len(tc.filename)-len(filepath.Ext(tc.filename))])
		mock, found := mockResponses[endpoint]
		if !found {
			t.Errorf("Expected mock response for endpoint %s not found", endpoint)
			continue
		}

		if mock.Method != tc.expected.Method {
			t.Errorf("For endpoint %s: expected method %s, got %s", endpoint, tc.expected.Method, mock.Method)
		}

		if len(mock.Responses) != len(tc.expected.Responses) {
			t.Errorf("For endpoint %s: expected %d responses, got %d", endpoint, len(tc.expected.Responses), len(mock.Responses))
			continue
		}

		for i, resp := range mock.Responses {
			expected := tc.expected.Responses[i]
			if resp.Status != expected.Status {
				t.Errorf("For endpoint %s, response %d: expected status %d, got %d", endpoint, i, expected.Status, resp.Status)
			}
			if resp.Description != expected.Description {
				t.Errorf("For endpoint %s, response %d: expected description %s, got %s", endpoint, i, expected.Description, resp.Description)
			}
		}
	}

	// Test with non-existent directory
	_, err = LoadResponses("/non-existent-directory")
	if err == nil {
		t.Error("Expected error when loading from non-existent directory, got nil")
	}
}

func TestFindResponse(t *testing.T) {
	// Create a test response with multiple configurations
	response := Response{
		Method: "POST",
		Responses: []ResponseConfig{
			{
				Status:      201,
				Body:        map[string]interface{}{"message": "Success"},
				InputBody:   map[string]interface{}{"name": "Test User", "email": "test@example.com"},
				Description: "Success case",
			},
			{
				Status:      400,
				Body:        map[string]interface{}{"error": "Invalid input"},
				InputBody:   map[string]interface{}{"name": "Test User"},
				Description: "Missing email",
			},
			{
				Status:      201,
				Body:        map[string]interface{}{"message": "Default success"},
				Description: "Default case",
			},
		},
	}

	testCases := []struct {
		name           string
		inputBody      interface{}
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:           "Matching input body",
			inputBody:      map[string]interface{}{"name": "Test User", "email": "test@example.com"},
			expectedStatus: 201,
			expectedBody:   map[string]interface{}{"message": "Success"},
		},
		{
			name:           "Partial matching input body",
			inputBody:      map[string]interface{}{"name": "Test User"},
			expectedStatus: 400,
			expectedBody:   map[string]interface{}{"error": "Invalid input"},
		},
		{
			name:           "No matching input body",
			inputBody:      map[string]interface{}{"email": "test@example.com"},
			expectedStatus: 201,
			expectedBody:   map[string]interface{}{"message": "Default success"},
		},
		{
			name:           "Nil input body",
			inputBody:      nil,
			expectedStatus: 201,
			expectedBody:   map[string]interface{}{"message": "Success"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := response.FindResponse(tc.inputBody)
			if result == nil {
				t.Fatal("Expected non-nil response")
			}

			if result.Status != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, result.Status)
			}

			// Compare response bodies
			expectedJSON, _ := json.Marshal(tc.expectedBody)
			actualJSON, _ := json.Marshal(result.Body)
			if string(expectedJSON) != string(actualJSON) {
				t.Errorf("Expected body %v, got %v", tc.expectedBody, result.Body)
			}
		})
	}
}
