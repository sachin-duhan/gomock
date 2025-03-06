package server

// EndpointInfo represents the structure of endpoint information
type EndpointInfo struct {
	Method    string         `json:"method"`
	Responses []ResponseInfo `json:"responses"`
}

// ResponseInfo represents the structure of response information
type ResponseInfo struct {
	Status    int         `json:"status"`
	InputBody interface{} `json:"input_body,omitempty"`
	Body      interface{} `json:"response_body,omitempty"`
}

// EndpointsResponse represents the response structure for the /endpoints route
type EndpointsResponse struct {
	Status    string                  `json:"status"`
	Endpoints map[string]EndpointInfo `json:"endpoints"`
}
