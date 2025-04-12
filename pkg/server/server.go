package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/sachin-duhan/gomock/pkg/mock"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Server represents the mock server
type Server struct {
	responses map[string]mock.Response
	port      string
	logger    *zap.Logger
	server    *http.Server
}

// New creates a new mock server instance
func New(responses map[string]mock.Response, port string) (*Server, error) {
	logger, err := initLogger()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %v", err)
	}

	return &Server{
		responses: responses,
		port:      port,
		logger:    logger,
	}, nil
}

// Start starts the mock server
func (s *Server) Start() error {
	// Set up server routes
	mux := http.NewServeMux()
	mux.HandleFunc("/endpoints", s.handleEndpointsList)
	mux.HandleFunc("/", s.handleMockRequest)

	// Create HTTP server
	s.server = &http.Server{
		Addr:    ":" + s.port,
		Handler: s.logMiddleware(mux),
	}

	// Start the server
	s.logger.Info("Starting mock server", zap.String("port", s.port))
	return s.server.ListenAndServe()
}

// Stop gracefully shuts down the server
func (s *Server) Stop(ctx context.Context) error {
	if s.server != nil {
		s.logger.Info("Shutting down server")
		return s.server.Shutdown(ctx)
	}
	return nil
}

// initLogger initializes the zap logger based on environment variables
func initLogger() (*zap.Logger, error) {
	// Get log path from environment variable, default to "logs" directory
	logPath := os.Getenv("LOG_PATH")
	if logPath == "" {
		logPath = "logs"
	}

	// Create logs directory if it doesn't exist
	if err := os.MkdirAll(logPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %v", err)
	}

	// Configure encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Create file core
	logFile := filepath.Join(logPath, "mock-server.log")
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %v", err)
	}

	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(file),
		zap.InfoLevel,
	)

	// Create console core for development
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		zap.InfoLevel,
	)

	// Combine cores
	core := zapcore.NewTee(fileCore, consoleCore)

	// Create logger
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	return logger, nil
}

// logMiddleware logs incoming requests and their responses
func (s *Server) logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a response wrapper to capture the status code
		rw := newResponseWriter(w)

		// Log request details
		s.logger.Info("Incoming request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("remote_addr", r.RemoteAddr),
			zap.String("user_agent", r.UserAgent()),
		)

		// Process the request
		start := time.Now()
		next.ServeHTTP(rw, r)
		duration := time.Since(start)

		// Log response details
		s.logger.Info("Request completed",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Int("status", rw.status),
			zap.Duration("duration", duration),
		)
	})
}

// responseWriter wraps http.ResponseWriter to capture the status code
type responseWriter struct {
	http.ResponseWriter
	status int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
