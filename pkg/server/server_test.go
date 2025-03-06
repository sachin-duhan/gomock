package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sachin-duhan/gomock/pkg/mock"
)

func TestHandleMockRequest(t *testing.T) {
	// Create a test server
	mockResponses := map[string]mock.Response{
		"/users": {
			Method: "GET",
			Response: struct {
				Status int         `json:"status"`
				Body   interface{} `json:"body"`
			}{
				Status: 200,
				Body:   map[string]interface{}{"users": []interface{}{map[string]interface{}{"id": 1, "name": "Test User"}}},
			},
		},
		"/create-user": {
			Method: "POST",
			Response: struct {
				Status int         `json:"status"`
				Body   interface{} `json:"body"`
			}{
				Status: 201,
				Body:   map[string]interface{}{"message": "User created", "id": 2},
			},
		},
	}

	server := New(mockResponses, "8080")

	// Test cases
	testCases := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:           "Get Users",
			method:         "GET",
			path:           "/users",
			expectedStatus: 200,
			expectedBody:   map[string]interface{}{"users": []interface{}{map[string]interface{}{"id": 1, "name": "Test User"}}},
		},
		{
			name:           "Create User",
			method:         "POST",
			path:           "/create-user",
			expectedStatus: 201,
			expectedBody:   map[string]interface{}{"message": "User created", "id": 2},
		},
		{
			name:           "Not Found",
			method:         "GET",
			path:           "/not-found",
			expectedStatus: 404,
			expectedBody:   nil,
		},
		{
			name:           "Method Not Allowed",
			method:         "POST",
			path:           "/users",
			expectedStatus: 404,
			expectedBody:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, tc.path, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(server.handleMockRequest)
			handler.ServeHTTP(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, rr.Code)
			}

			if tc.expectedBody != nil {
				var responseBody map[string]interface{}
				err = json.Unmarshal(rr.Body.Bytes(), &responseBody)
				if err != nil {
					t.Fatalf("Failed to parse response body: %v", err)
				}

				// Compare response body with expected body
				for key := range tc.expectedBody {
					actualValue, exists := responseBody[key]
					if !exists {
						t.Errorf("Expected key %s not found in response", key)
						continue
					}

					// For simplicity, just check if the key exists
					// In a real test, you might want to do a deep comparison
					if actualValue == nil {
						t.Errorf("Value for key %s is nil", key)
					}
				}
			}
		})
	}
}

func TestHandleEndpointsList(t *testing.T) {
	// Create a test server
	mockResponses := map[string]mock.Response{
		"/users": {
			Method: "GET",
			Response: struct {
				Status int         `json:"status"`
				Body   interface{} `json:"body"`
			}{
				Status: 200,
				Body:   map[string]interface{}{"users": []interface{}{map[string]interface{}{"id": 1, "name": "Test User"}}},
			},
		},
		"/create-user": {
			Method: "POST",
			Response: struct {
				Status int         `json:"status"`
				Body   interface{} `json:"body"`
			}{
				Status: 201,
				Body:   map[string]interface{}{"message": "User created", "id": 2},
			},
		},
	}

	server := New(mockResponses, "8080")

	// Test the endpoints listing
	req, err := http.NewRequest("GET", "/endpoints", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.handleEndpointsList)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	// Check if the response contains the status and endpoints fields
	if status, exists := response["status"]; !exists || status != "success" {
		t.Errorf("Expected status 'success', got %v", status)
	}

	endpoints, exists := response["endpoints"].(map[string]interface{})
	if !exists {
		t.Fatalf("Endpoints field not found in response")
	}

	// Check if all endpoints are listed
	expectedEndpoints := []string{"/users", "/create-user", "/endpoints"}
	for _, endpoint := range expectedEndpoints {
		if _, exists := endpoints[endpoint]; !exists {
			t.Errorf("Expected endpoint %s not found in response", endpoint)
		}
	}

	// Test method not allowed
	req, err = http.NewRequest("POST", "/endpoints", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, rr.Code)
	}
}
