package mock

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// Response represents a mock API response configuration
type Response struct {
	Method    string           `json:"method"`
	Responses []ResponseConfig `json:"responses"`
}

// ResponseConfig represents a specific response configuration for an endpoint
type ResponseConfig struct {
	Status      int         `json:"status"`
	Body        interface{} `json:"body"`
	InputBody   interface{} `json:"input_body,omitempty"`
	Description string      `json:"description,omitempty"`
}

// LoadResponses loads mock responses from JSON files in the specified directory
func LoadResponses(path string) (map[string]Response, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	mockResponses := make(map[string]Response)

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			filePath := filepath.Join(path, file.Name())
			content, err := ioutil.ReadFile(filePath)
			if err != nil {
				return nil, err
			}

			var mock Response
			err = json.Unmarshal(content, &mock)
			if err != nil {
				return nil, err
			}

			// Extract the endpoint path from the filename (without extension)
			endpoint := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
			// If the endpoint doesn't start with /, add it
			if !strings.HasPrefix(endpoint, "/") {
				endpoint = "/" + endpoint
			}
			mockResponses[endpoint] = mock
		}
	}

	return mockResponses, nil
}

// FindResponse finds the appropriate response based on input body
func (r *Response) FindResponse(inputBody interface{}) *ResponseConfig {
	// If no responses defined, return nil
	if len(r.Responses) == 0 {
		return nil
	}

	// For GET requests or no input body, return the first response
	if r.Method == "GET" || inputBody == nil {
		return &r.Responses[0]
	}

	// Try to find a response with matching input body
	for _, resp := range r.Responses {
		if resp.InputBody != nil {
			// Convert both to JSON for comparison
			inputJSON, err := json.Marshal(inputBody)
			if err != nil {
				continue
			}
			configJSON, err := json.Marshal(resp.InputBody)
			if err != nil {
				continue
			}
			if string(inputJSON) == string(configJSON) {
				return &resp
			}
		}
	}

	// If no matching input body found, return the last response as default
	return &r.Responses[len(r.Responses)-1]
}
