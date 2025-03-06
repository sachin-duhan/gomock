# Gomock

## Project Overview
This project implements a simple mock server in Go that reads mock responses from JSON files stored in a specified folder. The server responds to HTTP requests based on paths defined in the JSON files, making it a great tool for testing and simulating APIs during development. The folder path for the JSON files is configurable via environment variables.

## Features
- **API Mocking**: Simulate API responses for testing and development purposes
- **Response Simulation**: Read predefined responses from JSON files
- **Environment-specific Configuration**: Configure folder paths and port through environment variables
- **Endpoints Listing**: Built-in `/endpoints` route to list all available endpoints
- **Modular Design**: Well-structured code for easy maintenance and extension
- **Automated Testing**: Comprehensive test suite to ensure reliability
- **CI/CD Integration**: GitHub Actions workflow for automated testing and building

## Project Structure
The project follows a modular structure:
```
gomock/
├── main.go                 # Entry point
├── pkg/
│   ├── config/             # Configuration package
│   │   ├── config.go
│   │   └── config_test.go
│   ├── mock/               # Mock response handling
│   │   ├── mock.go
│   │   └── mock_test.go
│   └── server/             # HTTP server
│       ├── server.go
│       └── server_test.go
├── endpoints/              # Mock response JSON files
│   ├── users.json
│   ├── user-details.json
│   └── create-user.json
└── ...
```

## Use Case
- **API Mocking**: This mock server can be used to simulate API responses for testing and development purposes.
- **Response Simulation**: The server reads predefined responses from JSON files, making it easy to simulate different API scenarios for various endpoints.
- **Environment-specific Configuration**: The folder path for the JSON files can be defined through an environment variable, allowing for flexible configuration across different environments.

## Configuration
The server can be configured using environment variables:
- `JSON_FOLDER_PATH`: Path to the folder containing JSON files (default: `./endpoints`)
- `PORT`: Port on which the server will listen (default: `8080`)

You can set these variables in a `.env` file or directly in your environment.

## JSON File Structure
Each JSON file in the specified folder represents a single endpoint. The filename (without extension) is used as the endpoint path. For example, a file named `users.json` will be mapped to the `/users` endpoint.

### JSON File Format
```json
{
  "method": "POST",
  "responses": [
    {
      "status": 201,
      "body": {
        "message": "User created successfully",
        "user_id": 3
      },
      "input_body": {
        "name": "John Doe",
        "email": "john@example.com",
        "role": "user"
      },
      "description": "Successfully creates a new user with valid input"
    },
    {
      "status": 400,
      "body": {
        "error": "Invalid input",
        "details": [
          "Email is required",
          "Name is required"
        ]
      },
      "input_body": {
        "role": "user"
      },
      "description": "Returns error when required fields are missing"
    }
  ]
}
```

- `method`: HTTP method (GET, POST, PUT, DELETE, etc.)
- `responses`: Array of possible responses for the endpoint
  - `status`: HTTP status code
  - `body`: Response body (can be any valid JSON)
  - `input_body`: (Optional) Expected request body to match this response
  - `description`: (Optional) Description of when this response is returned

### Response Selection
The server selects the appropriate response based on:
1. If the request has a body, it tries to match it with the `input_body` of responses
2. If a match is found, that response is returned
3. If no match is found or no body is provided, the first response in the array is returned

### Sample JSON Files

1. **Example 1 - `users.json`**:
```json
{
  "method": "GET",
  "responses": [
    {
      "status": 200,
      "body": {
        "users": [
          {
            "id": 1,
            "name": "John Doe",
            "email": "john@example.com",
            "role": "admin"
          },
          {
            "id": 2,
            "name": "Jane Smith",
            "email": "jane@example.com",
            "role": "user"
          }
        ]
      },
      "description": "Returns list of users when authenticated"
    },
    {
      "status": 401,
      "body": {
        "error": "Unauthorized",
        "message": "Authentication required"
      },
      "description": "Returns error when not authenticated"
    }
  ]
}
```

2. **Example 2 - `create-user.json`**:
```json
{
  "method": "POST",
  "responses": [
    {
      "status": 201,
      "body": {
        "message": "User created successfully",
        "user_id": 3
      },
      "input_body": {
        "name": "John Doe",
        "email": "john@example.com",
        "role": "user"
      },
      "description": "Successfully creates a new user with valid input"
    },
    {
      "status": 400,
      "body": {
        "error": "Invalid input",
        "details": [
          "Email is required",
          "Name is required"
        ]
      },
      "input_body": {
        "role": "user"
      },
      "description": "Returns error when required fields are missing"
    }
  ]
}
```

## Endpoints
The server provides the following built-in endpoints:
- `/endpoints`: Lists all available endpoints with their methods, status codes, and sample curl commands

## Example Endpoints Response
```json
{
  "status": "success",
  "endpoints": {
    "/users": {
      "method": "GET",
      "status_code": 200,
      "sample_curl": "curl -X GET http://localhost:8080/users"
    },
    "/create-user": {
      "method": "POST",
      "status_code": 201,
      "sample_curl": "curl -X POST -H \"Content-Type: application/json\" -d '{\"key\":\"value\"}' http://localhost:8080/create-user"
    },
    "/endpoints": {
      "method": "GET",
      "status_code": 200,
      "description": "Lists all available endpoints with their methods, status codes, and sample curl commands",
      "sample_curl": "curl -X GET http://localhost:8080/endpoints"
    }
  }
}
```

## Setup Commands

1. Build the Docker image:
   ```bash
   docker build -t mock-server .
   ```

2. Run the Docker container:
   ```bash
   docker run -p 8080:8080 --env-file .env mock-server
   ```

## Running Locally
1. Clone the repository
2. Create a `.env` file with your configuration (or use the defaults)
3. Run the server:
   ```bash
   go run main.go
   ```

## Example Usage
1. Start the server
2. Access the endpoints listing: `http://localhost:8080/endpoints`
3. Access a mock endpoint: `http://localhost:8080/users`

## Testing
The project includes a comprehensive test suite to ensure reliability. To run the tests:

```bash
go test -v ./...
```

The tests cover:
- **Configuration**: Testing environment variable loading and defaults
- **Mock Responses**: Testing loading and parsing of mock response files
- **HTTP Server**: Testing request handling and response generation
- **Integration**: Testing the integration of all components

## CI/CD
The project includes a GitHub Actions workflow that automatically runs tests and builds the application on push to the master branch. The workflow is defined in `.github/workflows/go.yml` and includes:

- Running tests
- Running linter
- Building the application 