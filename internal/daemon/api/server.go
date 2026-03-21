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

type Server struct {
	db       *sql.DB
	listener net.Listener
	mux      *http.ServeMux
}

func NewServer(db *sql.DB, socketPath string) (*Server, error) {
	socketDir := filepath.Dir(socketPath)
	if err := os.MkdirAll(socketDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create socket directory: %w", err)
	}

	os.Remove(socketPath)

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create Unix socket: %w", err)
	}

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

func (s *Server) registerRoutes() {
	s.mux.HandleFunc("/api/v1/health", s.handleHealth)
	s.mux.HandleFunc("/api/v1/stats/today", s.handleGetTodayStats)
	s.mux.HandleFunc("/api/v1/stats/daily/", s.handleGetDailyStatsRoute)
}

func (s *Server) Start() error {
	logger.Verbose("API server starting...")
	return http.Serve(s.listener, s.mux)
}

func (s *Server) Close() error {
	logger.Verbose("API server closing...")
	return s.listener.Close()
}
