package models

// AppData represents application usage data for the frontend
type AppData struct {
	Class              string  `json:"class"`
	TotalTime          int64   `json:"total_time"`
	TotalTimeFormatted string  `json:"total_time_formatted"`
	OpenCount          int64   `json:"open_count"`
	LastSeen           string  `json:"last_seen"`
	Percentage         float64 `json:"percentage"`
}

// DailyData represents daily usage data
type DailyData struct {
	Date               string    `json:"date"`
	TotalTime          int64     `json:"total_time"`
	TotalTimeFormatted string    `json:"total_time_formatted"`
	Apps               []AppData `json:"apps"`
}
