package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"hyprtime/internal/database"
	"hyprtime/internal/models"
)

type ScreenTimeService struct {
	db *sql.DB
}

// NewScreenTimeService creates a new screen time service
func NewScreenTimeService() *ScreenTimeService {
	db, err := database.InitDB()
	if err != nil {
		log.Printf("Warning: Failed to initialize database: %v", err)
		return &ScreenTimeService{db: nil}
	}
	return &ScreenTimeService{db: db}
}

// Close closes the database connection
func (s *ScreenTimeService) Close() {
	if s.db != nil {
		s.db.Close()
	}
}

// formatDuration formats seconds into a readable string
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

// GetOverview returns the overall screen time overview
func (s *ScreenTimeService) GetOverview() (*models.ScreenTimeOverview, error) {
	if s.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	// Get total screen time
	totalTime, err := database.GetTotalScreenTime(s.db)
	if err != nil {
		return nil, err
	}

	// Get today's screen time
	today := time.Now().Format("2006-01-02")
	todayTime, err := database.GetTotalScreenTimeForDate(s.db, today)
	if err != nil {
		return nil, err
	}

	// Get top apps
	apps, err := database.GetAllApps(s.db)
	if err != nil {
		return nil, err
	}

	topApps := make([]models.AppData, 0)
	for i, app := range apps {
		if i >= 10 { // Limit to top 10
			break
		}

		percentage := 0.0
		if totalTime > 0 {
			percentage = float64(app.TotalTime) / float64(totalTime) * 100
		}

		topApps = append(topApps, models.AppData{
			Class:              app.Class,
			TotalTime:          app.TotalTime,
			TotalTimeFormatted: formatDuration(app.TotalTime),
			OpenCount:          app.OpenCount,
			LastSeen:           app.LastSeen,
			Percentage:         percentage,
		})
	}

	return &models.ScreenTimeOverview{
		TotalTime:          totalTime,
		TotalTimeFormatted: formatDuration(totalTime),
		TodayTime:          todayTime,
		TodayTimeFormatted: formatDuration(todayTime),
		TopApps:            topApps,
	}, nil
}

// GetDailyStats returns statistics for a specific date
func (s *ScreenTimeService) GetDailyStats(date string) (*models.DailyData, error) {
	if s.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	stats, err := database.GetDailyStats(s.db, date)
	if err != nil {
		return nil, err
	}

	totalTime, err := database.GetTotalScreenTimeForDate(s.db, date)
	if err != nil {
		return nil, err
	}

	apps := make([]models.AppData, 0)
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
			Percentage:         percentage,
		})
	}

	return &models.DailyData{
		Date:               date,
		TotalTime:          totalTime,
		TotalTimeFormatted: formatDuration(totalTime),
		Apps:               apps,
	}, nil
}

// GetTodayStats returns today's statistics
func (s *ScreenTimeService) GetTodayStats() (*models.DailyData, error) {
	today := time.Now().Format("2006-01-02")
	return s.GetDailyStats(today)
}
