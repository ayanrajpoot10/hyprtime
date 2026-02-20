package utils

import (
	"os"
	"os/user"
	"path/filepath"
)

func GetUser() string {
	currentUser, err := user.Current()
	if err != nil {
		return os.Getenv("USER")
	}
	return currentUser.Username
}

func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

func GetDBPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = filepath.Join("/home", GetUser())
	}
	return filepath.Join(homeDir, ".local", "share", "hyprtime", "hyprtime.db")
}
