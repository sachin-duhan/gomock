package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sachin-duhan/gomock/pkg/mock"
	"go.uber.org/zap/zaptest"
)

func setupTestServer(t *testing.T) *Server {
	// Create a test logger
	logger := zaptest.NewLogger(t)

	mockResponses := map[string]mock.Response{
		"/users": {
			Method: "GET",
			Responses: []mock.ResponseConfig{
				{
					Status:      200,
					Body:        map[string]interface{}{"users": []interface{}{map[string]interface{}{"id": 1, "name": "Test User"}}},
					Description: "Success response",
				},
				{
					Status:      401,
					Body:        map[string]interface{}{"error": "Unauthorized"},
					Description: "Error response",
				},
				{
					Status:      403,
					Body:        map[string]interface{}{"error": "Forbidden"},
					Description: "Forbidden response",
				},
			},
		},
		"/create-user": {
			Method: "POST",
			Responses: []mock.ResponseConfig{
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
					Status:      409,
					Body:        map[string]interface{}{"error": "User already exists"},
					Description: "Conflict case",
				},
			},
		},
		"/products": {
			Method: "GET",
			Responses: []mock.ResponseConfig{
				{
					Status:      200,
					Body:        map[string]interface{}{"products": []interface{}{map[string]interface{}{"id": 1, "name": "Product 1"}}},
					Description: "Success response",
				},
			},
		},
	}

	return &Server{
		responses: mockResponses,
		port:      "8080",
		logger:    logger,
	}
}

func TestHandleMockRequest(t *testing.T) {
	server := setupTestServer(t)

	// Test cases
	testCases := []struct {
		name           string
		method         string
		path           string
		body           interface{}
		headers        map[string]string
		expectedStatus int
		expectedBody   interface{}
	}{
		// Basic functionality tests
		{
			name:           "Get Users Success - Default Response",
			method:         "GET",
			path:           "/users",
			expectedStatus: 200,
			expectedBody:   map[string]interface{}{"users": []interface{}{map[string]interface{}{"id": 1, "name": "Test User"}}},
		},
		{
			name:           "Create User Success - Default Response",
			method:         "POST",
			path:           "/create-user",
			body:           map[string]interface{}{"name": "Test User", "email": "test@example.com"},
			expectedStatus: 201,
			expectedBody:   map[string]interface{}{"message": "Success"},
		},

		// x-stub-status header tests - Existing status codes
		{
			name:           "Get Users - Force 401 via Header",
			method:         "GET",
			path:           "/users",
			headers:        map[string]string{"x-stub-status": "401"},
			expectedStatus: 401,
			expectedBody:   map[string]interface{}{"error": "Unauthorized"},
		},
		{
			name:           "Get Users - Force 403 via Header",
			method:         "GET",
			path:           "/users",
			headers:        map[string]string{"x-stub-status": "403"},
			expectedStatus: 403,
			expectedBody:   map[string]interface{}{"error": "Forbidden"},
		},
		{
			name:           "Create User - Force 409 via Header",
			method:         "POST",
			path:           "/create-user",
			body:           map[string]interface{}{"name": "Test User", "email": "test@example.com"},
			headers:        map[string]string{"x-stub-status": "409"},
			expectedStatus: 409,
			expectedBody:   map[string]interface{}{"error": "User already exists"},
		},

		// x-stub-status header tests - Non-existent status codes
		{
			name:           "Get Users - Non-existent Status Code (404)",
			method:         "GET",
			path:           "/users",
			headers:        map[string]string{"x-stub-status": "404"},
			expectedStatus: 200, // Should fall back to first response
			expectedBody:   map[string]interface{}{"users": []interface{}{map[string]interface{}{"id": 1, "name": "Test User"}}},
		},
		{
			name:           "Products - Non-existent Status Code (500)",
			method:         "GET",
			path:           "/products",
			headers:        map[string]string{"x-stub-status": "500"},
			expectedStatus: 200, // Should fall back to first response
			expectedBody:   map[string]interface{}{"products": []interface{}{map[string]interface{}{"id": 1, "name": "Product 1"}}},
		},

		// x-stub-status header tests - Invalid values
		{
			name:           "Get Users - Invalid Status Code (non-numeric)",
			method:         "GET",
			path:           "/users",
			headers:        map[string]string{"x-stub-status": "invalid"},
			expectedStatus: 200, // Should ignore invalid header and use default
			expectedBody:   map[string]interface{}{"users": []interface{}{map[string]interface{}{"id": 1, "name": "Test User"}}},
		},
		{
			name:           "Get Users - Empty Status Code",
			method:         "GET",
			path:           "/users",
			headers:        map[string]string{"x-stub-status": ""},
			expectedStatus: 200, // Should ignore empty header and use default
			expectedBody:   map[string]interface{}{"users": []interface{}{map[string]interface{}{"id": 1, "name": "Test User"}}},
		},

		// Error cases
		{
			name:           "Not Found Endpoint",
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
		{
			name:           "Invalid JSON Body",
			method:         "POST",
			path:           "/create-user",
			body:           "invalid json",
			expectedStatus: 400,
			expectedBody:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var req *http.Request
			var err error

			if tc.body != nil {
				var bodyBytes []byte
				if str, ok := tc.body.(string); ok {
					bodyBytes = []byte(str)
				} else {
					bodyBytes, _ = json.Marshal(tc.body)
				}
				req, err = http.NewRequest(tc.method, tc.path, bytes.NewBuffer(bodyBytes))
			} else {
				req, err = http.NewRequest(tc.method, tc.path, nil)
			}

			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			// Add custom headers if specified
			if tc.headers != nil {
				for key, value := range tc.headers {
					req.Header.Set(key, value)
				}
			}

			// Set Content-Type for POST requests
			if tc.method == "POST" && tc.body != nil {
				req.Header.Set("Content-Type", "application/json")
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
				expectedJSON, _ := json.Marshal(tc.expectedBody)
				actualJSON, _ := json.Marshal(responseBody)
				if string(expectedJSON) != string(actualJSON) {
					t.Errorf("Expected body %v, got %v", tc.expectedBody, responseBody)
				}
			}
		})
	}
}

func TestHandleEndpointsList(t *testing.T) {
	server := setupTestServer(t)

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

	var response EndpointsResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	// Check if the response contains the status and endpoints fields
	if response.Status != "success" {
		t.Errorf("Expected status 'success', got %v", response.Status)
	}

	// Check if all endpoints are listed
	expectedEndpoints := []string{"/users", "/create-user", "/products", "/endpoints"}
	for _, endpoint := range expectedEndpoints {
		if _, exists := response.Endpoints[endpoint]; !exists {
			t.Errorf("Expected endpoint %s not found in response", endpoint)
		}
	}

	// Verify endpoint structure
	for path, endpoint := range response.Endpoints {
		if endpoint.Method == "" {
			t.Errorf("Endpoint %s: missing method", path)
		}
		if len(endpoint.Responses) == 0 {
			t.Errorf("Endpoint %s: no responses defined", path)
		}
		for i, resp := range endpoint.Responses {
			if resp.Status == 0 {
				t.Errorf("Endpoint %s, response %d: missing status code", path, i)
			}
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
