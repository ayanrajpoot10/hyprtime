package database

import "database/sql"

// DailyStats holds per-app usage data for a single day.
type DailyStats struct {
	ID        int64  `json:"id"`
	AppID     int64  `json:"app_id"`
	Class     string `json:"class"`
	Date      string `json:"date"`
	TotalTime int64  `json:"total_time"`
	OpenCount int64  `json:"open_count"`
}

// GetDailyStats returns all app stats for a given date, ordered by time descending.
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
		var s DailyStats
		if err := rows.Scan(&s.ID, &s.AppID, &s.Class, &s.Date, &s.TotalTime, &s.OpenCount); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}

	return stats, rows.Err()
}

// GetTotalScreenTimeForDate returns the sum of all app time for a given date.
func GetTotalScreenTimeForDate(db *sql.DB, date string) (int64, error) {
	var total int64
	err := db.QueryRow(`
		SELECT COALESCE(SUM(total_time), 0)
		FROM daily_stats
		WHERE date = ?
	`, date).Scan(&total)
	return total, err
}
