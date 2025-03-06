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

## JSON File Structure

Each JSON file in your endpoints folder represents one endpoint. The filename becomes the endpoint path (e.g., `users.json` → `/users`).

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



## Development

1. Clone the repo
2. Install dependencies: `go mod download`
3. Run tests: `go test ./...`
4. Start server: `go run main.go`

### List Available Endpoints
```bash
curl http://localhost:8080/endpoints
```
