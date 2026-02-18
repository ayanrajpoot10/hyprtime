package database

import (
	"database/sql"
	"fmt"

	"hyprtime/internal/utils"

	_ "github.com/mattn/go-sqlite3"
)

// InitDB initializes the SQLite database connection
func InitDB() (*sql.DB, error) {
	user := utils.GetUser()
	dbFile := fmt.Sprintf("/home/%s/.local/share/hyprscreentime/screentime.db", user)

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// AppStats represents statistics for an application
type AppStats struct {
	ID        int64  `json:"id"`
	Class     string `json:"class"`
	TotalTime int64  `json:"total_time"`
	OpenCount int64  `json:"open_count"`
	LastSeen  string `json:"last_seen"`
	CreatedAt string `json:"created_at"`
}

// DailyStats represents daily statistics for an application
type DailyStats struct {
	ID        int64  `json:"id"`
	AppID     int64  `json:"app_id"`
	Class     string `json:"class"`
	Date      string `json:"date"`
	TotalTime int64  `json:"total_time"`
	OpenCount int64  `json:"open_count"`
}

// GetAllApps retrieves all app statistics
func GetAllApps(db *sql.DB) ([]AppStats, error) {
	rows, err := db.Query(`
		SELECT id, class, total_time, open_count, 
		       COALESCE(last_seen, ''), COALESCE(created_at, '')
		FROM apps
		ORDER BY total_time DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apps []AppStats
	for rows.Next() {
		var app AppStats
		err := rows.Scan(&app.ID, &app.Class, &app.TotalTime, &app.OpenCount, &app.LastSeen, &app.CreatedAt)
		if err != nil {
			return nil, err
		}
		apps = append(apps, app)
	}

	return apps, rows.Err()
}

// GetDailyStats retrieves daily statistics for a specific date
func GetDailyStats(db *sql.DB, date string) ([]DailyStats, error) {
	rows, err := db.Query(`
		SELECT ds.id, ds.app_id, a.class, ds.date, ds.total_time, ds.open_count
		FROM daily_stats ds
		JOIN apps a ON ds.app_id = a.id
		WHERE ds.date = ?
		ORDER BY ds.total_time DESC
	`, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []DailyStats
	for rows.Next() {
		var stat DailyStats
		err := rows.Scan(&stat.ID, &stat.AppID, &stat.Class, &stat.Date, &stat.TotalTime, &stat.OpenCount)
		if err != nil {
			return nil, err
		}
		stats = append(stats, stat)
	}

	return stats, rows.Err()
}

// GetTotalScreenTime returns the total screen time across all apps
func GetTotalScreenTime(db *sql.DB) (int64, error) {
	var total int64
	err := db.QueryRow("SELECT COALESCE(SUM(total_time), 0) FROM apps").Scan(&total)
	return total, err
}

// GetTotalScreenTimeForDate returns the total screen time for a specific date
func GetTotalScreenTimeForDate(db *sql.DB, date string) (int64, error) {
	var total int64
	err := db.QueryRow(`
		SELECT COALESCE(SUM(total_time), 0) 
		FROM daily_stats 
		WHERE date = ?
	`, date).Scan(&total)
	return total, err
}
