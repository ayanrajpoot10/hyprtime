package database

import (
	"database/sql"
	"fmt"
	"time"

	"hyprtime/internal/utils"

	_ "github.com/mattn/go-sqlite3"
)

// InitDB initializes the SQLite database and creates tables if needed
func InitDB() (*sql.DB, error) {
	user := utils.GetUser()
	dbFile := fmt.Sprintf("/home/%s/.local/share/hyprtime/hyprtime.db", user)

	// Create directory if it doesn't exist
	if err := utils.EnsureDir(fmt.Sprintf("/home/%s/.local/share/hyprtime", user)); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Create tables
	if err := createTables(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS apps (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		class TEXT NOT NULL UNIQUE,
		total_time INTEGER DEFAULT 0,
		open_count INTEGER DEFAULT 0,
		last_seen TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS daily_stats (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		app_id INTEGER NOT NULL,
		date DATE NOT NULL,
		total_time INTEGER DEFAULT 0,
		open_count INTEGER DEFAULT 0,
		FOREIGN KEY (app_id) REFERENCES apps(id),
		UNIQUE(app_id, date)
	);

	CREATE INDEX IF NOT EXISTS idx_daily_stats_date ON daily_stats(date);
	CREATE INDEX IF NOT EXISTS idx_daily_stats_app_date ON daily_stats(app_id, date);
	`

	_, err := db.Exec(schema)
	return err
}

// GetOrCreateApp gets or creates an app record
func GetOrCreateApp(db *sql.DB, class string) (int64, error) {
	var appID int64
	err := db.QueryRow("SELECT id FROM apps WHERE class = ?", class).Scan(&appID)
	if err == sql.ErrNoRows {
		result, err := db.Exec(
			"INSERT INTO apps (class, last_seen) VALUES (?, ?)",
			class, time.Now(),
		)
		if err != nil {
			return 0, err
		}
		return result.LastInsertId()
	}
	if err != nil {
		return 0, err
	}

	// Update last seen
	_, err = db.Exec("UPDATE apps SET last_seen = ? WHERE id = ?", time.Now(), appID)
	return appID, err
}

// UpdateAppTime updates the total time for an app and its daily stats
func UpdateAppTime(db *sql.DB, appID int64, duration int64) error {
	if duration <= 0 {
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		"UPDATE apps SET total_time = total_time + ?, last_seen = ? WHERE id = ?",
		duration, time.Now(), appID,
	)
	if err != nil {
		return err
	}

	today := time.Now().Format("2006-01-02")
	_, err = tx.Exec(`
		INSERT INTO daily_stats (app_id, date, total_time, open_count)
		VALUES (?, ?, ?, 0)
		ON CONFLICT(app_id, date) DO UPDATE SET
			total_time = total_time + ?
	`, appID, today, duration, duration)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// IncrementOpenCount increments the open count for an app
func IncrementOpenCount(db *sql.DB, appID int64) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("UPDATE apps SET open_count = open_count + 1 WHERE id = ?", appID)
	if err != nil {
		return err
	}

	today := time.Now().Format("2006-01-02")
	_, err = tx.Exec(`
		INSERT INTO daily_stats (app_id, date, total_time, open_count)
		VALUES (?, ?, 0, 1)
		ON CONFLICT(app_id, date) DO UPDATE SET
			open_count = open_count + 1
	`, appID, today)
	if err != nil {
		return err
	}

	return tx.Commit()
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
