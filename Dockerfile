FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/gomock

# Create a minimal production image
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/gomock .

# Create directory for mock responses
RUN mkdir -p /app/endpoints

# Set default environment variables
ENV PORT=8080
ENV JSON_FOLDER_PATH=/app/endpoints

# Expose the port
EXPOSE ${PORT}

# Command to run the application
CMD ["./gomock"]
