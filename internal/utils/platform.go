package utils

import (
	"os"
	"path/filepath"
	"runtime"
)

func GetConfigDir() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(os.Getenv("APPDATA"), "iwashere")
	}
	// Linux/macOS
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "iwashere")
}

// Disable colors on Windows if needed
func UseColors() bool {
	if runtime.GOOS == "windows" {
		// Check if we're in a modern terminal that supports colors
		return os.Getenv("TERM") != "" || os.Getenv("WT_SESSION") != ""
	}
	return true
}
