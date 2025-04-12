# Gomock - Simple HTTP Mock Server

A lightweight mock server that lets you quickly set up HTTP endpoints with configurable responses. Perfect for development and testing when you need to simulate API responses.

## Use Case
- **API Mocking**: This mock server can be used to simulate API responses for testing and development purposes.
- **Response Simulation**: The server reads predefined responses from JSON files, making it easy to simulate different API scenarios for various endpoints.
- **Environment-specific Configuration**: The folder path for the JSON files can be defined through an environment variable, allowing for flexible configuration across different environments.

## Setup
```bash
git clone github.com/sachin-duhan/gomock
cd gomock
go mod tidy
make test
make run
```

## Quick Start
1. Build:
```bash
docker build -t gomock .
```

2. Run:
```bash
# Using default settings
docker run -p 8080:8080 gomock

# With custom JSON folder
docker run -p 8080:8080 -v $(pwd)/yours-endpoints:/app/endpoints gomock
```
## Features

- **Simple JSON Configuration**: Define endpoints and their responses in JSON files
- **Multiple Responses**: Support different responses based on input body or status code
- **Status Code Override**: Use `x-stub-resStatus` header to force specific status codes
- **Endpoints Discovery**: Built-in `/endpoints` route lists all available endpoints
- **Custom Endpoint Paths**: Define explicit API paths in your JSON files

## JSON File Structure

Each JSON file in your endpoints folder represents one endpoint. There are two ways to define the endpoint path:

1. **Default**: The filename becomes the endpoint path (e.g., `users.json` → `/users`)
2. **Custom**: Use the `path` property to explicitly define the endpoint path (e.g., `"path": "/api/v1/users"`)

The `path` property is especially useful for complex API paths with multiple segments.

### Basic Example (GET endpoint)
```json
{
  "method": "GET",
  "responses": [
    {
      "status": 200,
      "body": {"message": "Success"}
    }
  ]
}
```

### Example with Custom Path
```json
{
  "method": "POST",
  "path": "/api/v1/auth/token",
  "responses": [
    {
      "status": 200,
      "body": {"token": "eyJ0eXAi..."}
    }
  ]
}
```

### Input Body Matching (POST endpoint)
```json
{
  "method": "POST",
  "responses": [
    {
      "status": 201,
      "input_body": {
        "name": "John",
        "email": "john@example.com"
      },
      "body": {
        "id": 1,
        "message": "User created"
      }
    },
    {
      "status": 400,
      "input_body": {
        "name": "John"
      },
      "body": {
        "error": "Email is required"
      }
    }
  ]
}
```

## Using the x-stub-resStatus Header

You can force a specific status code response by using the `x-stub-resStatus` header:

```bash
# Force a 401 response for the /users endpoint
curl -H "x-stub-resStatus: 401" http://localhost:8080/users
```

If the requested status code exists in the endpoint's responses, that response will be used. If not, it will fall back to the default response selection logic.

## Development

1. Clone the repo
2. Install dependencies: `go mod download`
3. Run tests: `go test ./...`
4. Start server: `go run main.go`

### List Available Endpoints
```bash
curl http://localhost:8080/endpoints
```
