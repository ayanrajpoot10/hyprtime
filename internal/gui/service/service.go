package service

import (
	"database/sql"
	"fmt"
	"time"

	"hyprtime/internal/shared/database"
	"hyprtime/internal/shared/models"
)

type ScreenTimeService struct {
	db *sql.DB
}

func NewScreenTimeService() *ScreenTimeService {
	db, err := database.InitDB()
	if err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
	}
	return &ScreenTimeService{
		db: db,
	}
}

func (s *ScreenTimeService) Close() {
	if s.db != nil {
		s.db.Close()
	}
}

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

func (s *ScreenTimeService) GetDailyStats(date string) (*models.DailyData, error) {
	if s.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	stats, err := database.GetDailyStats(s.db, date)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	totalTime, err := database.GetTotalScreenTimeForDate(s.db, date)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
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

	return &response, nil
}

func (s *ScreenTimeService) GetTodayStats() (*models.DailyData, error) {
	today := time.Now().Format("2006-01-02")
	return s.GetDailyStats(today)
}
