package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/sachin-duhan/gomock/pkg/mock"
	"go.uber.org/zap"
)

// handleMockRequest handles incoming API requests and returns mock responses
func (s *Server) handleMockRequest(w http.ResponseWriter, r *http.Request) {
	mock, err := s.findMockResponse(r.URL.Path)
	if err != nil {
		s.logger.Error("Mock response not found",
			zap.String("path", r.URL.Path),
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := s.validateMethod(r.Method, mock.Method); err != nil {
		s.logger.Error("Invalid HTTP method",
			zap.String("path", r.URL.Path),
			zap.String("expected_method", mock.Method),
			zap.String("actual_method", r.Method),
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	inputBody, err := s.parseRequestBody(r)
	if err != nil {
		s.logger.Error("Failed to parse request body",
			zap.String("path", r.URL.Path),
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.logger.Debug("Request body parsed",
		zap.String("path", r.URL.Path),
		zap.Any("input_body", inputBody),
	)

	// Get desired status code from header
	desiredStatus := 0
	if statusHeader := r.Header.Get("x-stub-status"); statusHeader != "" {
		if status, err := strconv.Atoi(statusHeader); err == nil {
			desiredStatus = status
			s.logger.Debug("Using status code from header",
				zap.Int("status", desiredStatus),
			)
		}
	}

	response := s.findMatchingResponse(mock, inputBody, desiredStatus)
	if response == nil {
		s.logger.Error("No matching response found",
			zap.String("path", r.URL.Path),
			zap.Any("input_body", inputBody),
		)
		http.Error(w, "No matching response found", http.StatusInternalServerError)
		return
	}

	s.logger.Debug("Found matching response",
		zap.String("path", r.URL.Path),
		zap.Int("status", response.Status),
		zap.Any("response_body", response.Body),
	)

	s.writeJSONResponse(w, response.Status, response.Body)
}

// handleEndpointsList returns a list of all available endpoints
func (s *Server) handleEndpointsList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.logger.Error("Invalid method for endpoints list",
			zap.String("method", r.Method),
		)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	endpoints := s.buildEndpointsList(r.Host)
	response := EndpointsResponse{
		Status:    "success",
		Endpoints: endpoints,
	}

	s.logger.Debug("Endpoints list generated",
		zap.Int("endpoint_count", len(endpoints)),
	)

	s.writeJSONResponse(w, http.StatusOK, response)
}

// Helper functions

func (s *Server) findMockResponse(path string) (*mock.Response, error) {
	mock, exists := s.responses[path]
	if !exists {
		return nil, fmt.Errorf("Not Found")
	}
	return &mock, nil
}

func (s *Server) validateMethod(requestMethod, mockMethod string) error {
	if requestMethod != mockMethod {
		return fmt.Errorf("Not Found")
	}
	return nil
}

func (s *Server) parseRequestBody(r *http.Request) (interface{}, error) {
	if r.Body == nil || r.ContentLength == 0 {
		return nil, nil
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("Invalid request body")
	}
	r.Body.Close()
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	if len(body) == 0 {
		return nil, nil
	}

	var inputBody interface{}
	if err := json.Unmarshal(body, &inputBody); err != nil {
		return nil, fmt.Errorf("Invalid JSON in request body")
	}

	return inputBody, nil
}

func (s *Server) findMatchingResponse(mock *mock.Response, inputBody interface{}, desiredStatus int) *mock.ResponseConfig {
	// If desired status is specified, try to find a response with that status first
	if desiredStatus > 0 {
		for _, resp := range mock.Responses {
			if resp.Status == desiredStatus {
				return &resp
			}
		}
		// If no response with desired status found, log a warning
		s.logger.Warn("No response found for desired status code",
			zap.Int("desired_status", desiredStatus),
		)
	}

	// Fall back to normal matching logic
	response := mock.FindResponse(inputBody)
	if response == nil && len(mock.Responses) > 0 {
		return &mock.Responses[len(mock.Responses)-1]
	}
	return response
}

func (s *Server) writeJSONResponse(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		s.logger.Error("Failed to encode response",
			zap.Error(err),
			zap.Any("body", body),
		)
	}
}

func (s *Server) buildEndpointsList(host string) map[string]EndpointInfo {
	endpoints := make(map[string]EndpointInfo)

	// Add mock endpoints
	for path, mock := range s.responses {
		endpoints[path] = s.buildEndpointInfo(path, mock)
	}

	// Add endpoints listing endpoint
	endpoints["/endpoints"] = EndpointInfo{
		Method: "GET",
		Responses: []ResponseInfo{
			{
				Status: http.StatusOK,
			},
		},
	}

	return endpoints
}

func (s *Server) buildEndpointInfo(path string, mock mock.Response) EndpointInfo {
	responses := make([]ResponseInfo, len(mock.Responses))
	for i, resp := range mock.Responses {
		responses[i] = ResponseInfo{
			Status:    resp.Status,
			InputBody: resp.InputBody,
			Body:      resp.Body,
		}
	}

	return EndpointInfo{
		Method:    mock.Method,
		Responses: responses,
	}
}
