package database

import (
	"database/sql"
	"fmt"

	"hyprtime/internal/utils"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() (*sql.DB, error) {
	user := utils.GetUser()
	dbFile := fmt.Sprintf("/home/%s/.local/share/hyprtime/hyprtime.db", user)

	if err := utils.EnsureDir(fmt.Sprintf("/home/%s/.local/share/hyprtime", user)); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

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
