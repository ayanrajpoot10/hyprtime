package api

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"hyprtime/internal/logger"
)

// Server represents the HTTP API server
type Server struct {
	db       *sql.DB
	listener net.Listener
	mux      *http.ServeMux
}

// NewServer creates a new API server with Unix socket listener
func NewServer(db *sql.DB, socketPath string) (*Server, error) {
	// Ensure socket directory exists
	socketDir := filepath.Dir(socketPath)
	if err := os.MkdirAll(socketDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create socket directory: %w", err)
	}

	// Remove old socket if exists
	os.Remove(socketPath)

	// Create Unix socket listener
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create Unix socket: %w", err)
	}

	// Set socket permissions (only user can access)
	if err := os.Chmod(socketPath, 0600); err != nil {
		listener.Close()
		return nil, fmt.Errorf("failed to set socket permissions: %w", err)
	}

	logger.Info("API server listening on: %s", socketPath)

	s := &Server{
		db:       db,
		listener: listener,
		mux:      http.NewServeMux(),
	}

	s.registerRoutes()
	return s, nil
}

// registerRoutes sets up all API routes
func (s *Server) registerRoutes() {
	s.mux.HandleFunc("/api/v1/health", s.handleHealth)
	s.mux.HandleFunc("/api/v1/stats/today", s.handleGetTodayStats)
	s.mux.HandleFunc("/api/v1/stats/daily/", s.handleGetDailyStatsRoute)
}

// Start begins serving HTTP requests
func (s *Server) Start() error {
	logger.Verbose("API server starting...")
	return http.Serve(s.listener, s.mux)
}

// Close gracefully shuts down the server
func (s *Server) Close() error {
	logger.Verbose("API server closing...")
	return s.listener.Close()
}
