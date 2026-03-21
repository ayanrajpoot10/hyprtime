package database

import (
	"database/sql"
	"time"
)

// GetOrCreateApp returns the ID of an existing app or inserts a new one.
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

	_, err = db.Exec("UPDATE apps SET last_seen = ? WHERE id = ?", time.Now(), appID)
	return appID, err
}

// UpdateAppTime adds duration (seconds) to an app's total and today's daily record.
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

// IncrementOpenCount bumps the open count for an app in both the global and daily tables.
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
