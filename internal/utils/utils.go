package utils

import (
	"os"
	"os/user"
)

// GetUser returns the current username
func GetUser() string {
	currentUser, err := user.Current()
	if err != nil {
		return "user"
	}
	return currentUser.Username
}

// EnsureDir creates a directory if it doesn't exist
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}
