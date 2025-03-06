# Use the official Golang image as the base image
FROM golang:1.20-alpine

# Set the working directory in the container
WORKDIR /app

# Copy the Go Modules and install dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the entire project
COPY . .

# Build the Go application
RUN go build -o mock-server .

# Expose port 8080 to access the server
EXPOSE 8080

# Run the server
CMD ["./mock-server"]
