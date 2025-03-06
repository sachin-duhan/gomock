package mock

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// Response represents a mock API response configuration
type Response struct {
	Method   string `json:"method"`
	Response struct {
		Status int         `json:"status"`
		Body   interface{} `json:"body"`
	} `json:"response"`
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
