package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"hyprtime/internal/gui/client"
	"hyprtime/internal/shared/models"
)

// ScreenTimeService provides screen time data to the GUI via API calls
type ScreenTimeService struct {
	client *client.Client
}

// getSocketPath returns the daemon's Unix socket path
func getSocketPath() string {
	runtimeDir := os.Getenv("XDG_RUNTIME_DIR")
	if runtimeDir == "" {
		runtimeDir = fmt.Sprintf("/run/user/%d", os.Getuid())
	}
	return filepath.Join(runtimeDir, "hyprtime", "daemon.sock")
}

// NewScreenTimeService creates a new service that communicates with the daemon via API
func NewScreenTimeService() *ScreenTimeService {
	socketPath := getSocketPath()
	return &ScreenTimeService{
		client: client.NewClient(socketPath),
	}
}

// Close cleans up the service (no-op for now, but kept for interface compatibility)
func (s *ScreenTimeService) Close() {
	// No cleanup needed for HTTP client
}

// GetDailyStats fetches daily stats for a specific date via daemon API
func (s *ScreenTimeService) GetDailyStats(date string) (*models.DailyData, error) {
	ctx := context.Background()
	data, err := s.client.GetDailyStats(ctx, date)
	if err != nil {
		return nil, fmt.Errorf("failed to get daily stats from daemon: %w", err)
	}
	return data, nil
}

// GetTodayStats fetches today's stats via daemon API
func (s *ScreenTimeService) GetTodayStats() (*models.DailyData, error) {
	ctx := context.Background()
	data, err := s.client.GetTodayStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get today's stats from daemon: %w", err)
	}
	return data, nil
}
