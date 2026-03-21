package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"hyprtime/internal/daemon/database"
	"hyprtime/internal/logger"
	"hyprtime/internal/shared/models"
)

func formatDuration(seconds int64) string {
	duration := time.Duration(seconds) * time.Second
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	secs := int(duration.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, secs)
	} else if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, secs)
	}
	return fmt.Sprintf("%ds", secs)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

func (s *Server) handleGetTodayStats(w http.ResponseWriter, r *http.Request) {
	today := time.Now().Format("2006-01-02")
	s.getDailyStats(w, today)
}

func (s *Server) handleGetDailyStatsRoute(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	prefix := "/api/v1/stats/daily/"

	if !strings.HasPrefix(path, prefix) {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	date := strings.TrimPrefix(path, prefix)
	if date == "" {
		http.Error(w, "Date is required", http.StatusBadRequest)
		return
	}

	s.getDailyStats(w, date)
}

func (s *Server) getDailyStats(w http.ResponseWriter, date string) {
	logger.Debug("API: GetDailyStats for date: %s", date)

	stats, err := database.GetDailyStats(s.db, date)
	if err != nil {
		logger.Error("Failed to get daily stats: %v", err)
		http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
		return
	}

	totalTime, err := database.GetTotalScreenTimeForDate(s.db, date)
	if err != nil {
		logger.Error("Failed to get total time: %v", err)
		http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
		return
	}

	apps := make([]models.AppData, 0, len(stats))
	for _, stat := range stats {
		percentage := 0.0
		if totalTime > 0 {
			percentage = float64(stat.TotalTime) / float64(totalTime) * 100
		}

		apps = append(apps, models.AppData{
			Class:              stat.Class,
			TotalTime:          stat.TotalTime,
			TotalTimeFormatted: formatDuration(stat.TotalTime),
			OpenCount:          stat.OpenCount,
			LastSeen:           stat.Date,
			Percentage:         percentage,
		})
	}

	response := models.DailyData{
		Date:               date,
		TotalTime:          totalTime,
		TotalTimeFormatted: formatDuration(totalTime),
		Apps:               apps,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	logger.Debug("API: Successfully returned stats for %s (%d apps)", date, len(apps))
}
