package utils

import (
	"os"
	"os/user"
	"path/filepath"
)

// GetUser returns the current username
func GetUser() string {
	currentUser, err := user.Current()
	if err != nil {
		// Fallback to USER env variable
		return os.Getenv("USER")
	}
	return currentUser.Username
}

// EnsureDir creates a directory if it doesn't exist
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// GetDBPath returns the database file path
func GetDBPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = filepath.Join("/home", GetUser())
	}
	return filepath.Join(homeDir, ".local", "share", "hyprtime", "hyprtime.db")
}
