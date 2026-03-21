package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"hyprtime/internal/gui/client"
	"hyprtime/internal/shared/models"
)

type ScreenTimeService struct {
	client *client.Client
}

func getSocketPath() string {
	runtimeDir := os.Getenv("XDG_RUNTIME_DIR")
	if runtimeDir == "" {
		runtimeDir = fmt.Sprintf("/run/user/%d", os.Getuid())
	}
	return filepath.Join(runtimeDir, "hyprtime", "daemon.sock")
}

func NewScreenTimeService() *ScreenTimeService {
	socketPath := getSocketPath()
	return &ScreenTimeService{
		client: client.NewClient(socketPath),
	}
}

func (s *ScreenTimeService) Close() {
}

func (s *ScreenTimeService) GetDailyStats(date string) (*models.DailyData, error) {
	ctx := context.Background()
	data, err := s.client.GetDailyStats(ctx, date)
	if err != nil {
		return nil, fmt.Errorf("failed to get daily stats from daemon: %w", err)
	}
	return data, nil
}

func (s *ScreenTimeService) GetTodayStats() (*models.DailyData, error) {
	ctx := context.Background()
	data, err := s.client.GetTodayStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get today's stats from daemon: %w", err)
	}
	return data, nil
}
