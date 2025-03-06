package mock

import (
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
				"response": {
					"status": 200,
					"body": {"users": [{"id": 1, "name": "Test User"}]}
				}
			}`,
			expected: Response{
				Method: "GET",
				Response: struct {
					Status int         `json:"status"`
					Body   interface{} `json:"body"`
				}{
					Status: 200,
					Body:   map[string]interface{}{"users": []interface{}{map[string]interface{}{"id": float64(1), "name": "Test User"}}},
				},
			},
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
			expected: Response{
				Method: "POST",
				Response: struct {
					Status int         `json:"status"`
					Body   interface{} `json:"body"`
				}{
					Status: 201,
					Body:   map[string]interface{}{"message": "User created", "id": float64(2)},
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

		if mock.Response.Status != tc.expected.Response.Status {
			t.Errorf("For endpoint %s: expected status %d, got %d", endpoint, tc.expected.Response.Status, mock.Response.Status)
		}
	}

	// Test with non-existent directory
	_, err = LoadResponses("/non-existent-directory")
	if err == nil {
		t.Error("Expected error when loading from non-existent directory, got nil")
	}
}
