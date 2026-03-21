package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"hyprtime/internal/shared/models"
)

// Client represents an HTTP client for communicating with the daemon API
type Client struct {
	httpClient *http.Client
	baseURL    string
	socketPath string
}

// NewClient creates a new API client that communicates over Unix socket
func NewClient(socketPath string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					return net.Dial("unix", socketPath)
				},
			},
		},
		baseURL:    "http://unix",
		socketPath: socketPath,
	}
}

// Health checks if the daemon API is responding
func (c *Client) Health(ctx context.Context) error {
	url := fmt.Sprintf("%s/api/v1/health", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to daemon at %s: %w", c.socketPath, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed with status: %d", resp.StatusCode)
	}

	return nil
}

// GetDailyStats fetches daily stats for a specific date (YYYY-MM-DD format)
func (c *Client) GetDailyStats(ctx context.Context, date string) (*models.DailyData, error) {
	url := fmt.Sprintf("%s/api/v1/stats/daily/%s", c.baseURL, date)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to daemon at %s: %w (is hyprtimed running?)", c.socketPath, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned error status: %d", resp.StatusCode)
	}

	var data models.DailyData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &data, nil
}

// GetTodayStats fetches stats for the current day
func (c *Client) GetTodayStats(ctx context.Context) (*models.DailyData, error) {
	url := fmt.Sprintf("%s/api/v1/stats/today", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to daemon at %s: %w (is hyprtimed running?)", c.socketPath, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned error status: %d", resp.StatusCode)
	}

	var data models.DailyData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &data, nil
}
