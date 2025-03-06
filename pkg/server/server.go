package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/sachin-duhan/gomock/pkg/mock"
)

// Server represents the mock server
type Server struct {
	mockResponses map[string]mock.Response
	port          string
}

// New creates a new mock server instance
func New(mockResponses map[string]mock.Response, port string) *Server {
	return &Server{
		mockResponses: mockResponses,
		port:          port,
	}
}

// handleMockRequest handles incoming API requests and returns mock responses
func (s *Server) handleMockRequest(w http.ResponseWriter, r *http.Request) {
	// Skip the endpoints listing route
	if r.URL.Path == "/endpoints" {
		return
	}

	// Find the mock response for the requested path
	mockResponse, found := s.mockResponses[r.URL.Path]
	if !found || mockResponse.Method != r.Method {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	// Set the response status code and write the response body
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(mockResponse.Response.Status)
	json.NewEncoder(w).Encode(mockResponse.Response.Body)
}

// handleEndpointsList returns a list of all available endpoints
func (s *Server) handleEndpointsList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Create a list of available endpoints with their methods, status codes, and sample curl commands
	type EndpointInfo struct {
		Method      string `json:"method"`
		StatusCode  int    `json:"status_code"`
		Description string `json:"description,omitempty"`
		SampleCurl  string `json:"sample_curl"`
	}

	// Get the host from the request or use localhost with the server's port
	host := r.Host
	if host == "" {
		host = "localhost:" + s.port
	}

	endpoints := make(map[string]EndpointInfo)
	for path, mock := range s.mockResponses {
		// Create a sample curl command based on the HTTP method
		var curlCmd string
		switch mock.Method {
		case http.MethodGet:
			curlCmd = fmt.Sprintf("curl -X GET http://%s%s", host, path)
		case http.MethodPost:
			curlCmd = fmt.Sprintf("curl -X POST -H \"Content-Type: application/json\" -d '{\"key\":\"value\"}' http://%s%s", host, path)
		case http.MethodPut:
			curlCmd = fmt.Sprintf("curl -X PUT -H \"Content-Type: application/json\" -d '{\"key\":\"value\"}' http://%s%s", host, path)
		case http.MethodDelete:
			curlCmd = fmt.Sprintf("curl -X DELETE http://%s%s", host, path)
		default:
			// Default to GET if method is not specified
			curlCmd = fmt.Sprintf("curl -X GET http://%s%s", host, path)
		}

		endpoints[path] = EndpointInfo{
			Method:     mock.Method,
			StatusCode: mock.Response.Status,
			SampleCurl: curlCmd,
		}
	}

	// Add the endpoints listing endpoint itself
	endpoints["/endpoints"] = EndpointInfo{
		Method:      "GET",
		StatusCode:  http.StatusOK,
		Description: "Lists all available endpoints with their methods, status codes, and sample curl commands",
		SampleCurl:  fmt.Sprintf("curl -X GET http://%s/endpoints", host),
	}

	// Return the list of endpoints as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "success",
		"endpoints": endpoints,
	})
}

// Start starts the mock server
func (s *Server) Start() error {
	// Set up server routes
	http.HandleFunc("/endpoints", s.handleEndpointsList)
	http.HandleFunc("/", s.handleMockRequest)

	// Start the server
	log.Printf("Starting mock server on port %s...", s.port)
	return http.ListenAndServe(":"+s.port, nil)
}
