//go:windows
//+build windows

package utils

import (
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/sys/windows"
)

// HideFile makes a file/folder hidden on supported platforms
func HideFile(path string) error {
	if runtime.GOOS == "windows" {
		// Windows: set hidden attribute
		// Convert to Windows absolute path
		absPath, err := filepath.Abs(path)
		if err != nil {
			return err
		}

		// Get file handle
		ptr, err := windows.UTF16PtrFromString(absPath)
		if err != nil {
			return err
		}

		// Set FILE_ATTRIBUTE_HIDDEN (0x2)
		return windows.SetFileAttributes(ptr, 0x2)
	}

	// On Unix, dot prefix is enough
	return nil
}

// IsHidden checks if a file/folder is hidden
func IsHidden(path string) bool {
	if runtime.GOOS == "windows" {
		// Windows: check hidden attribute
		absPath, _ := filepath.Abs(path)
		ptr, _ := windows.UTF16PtrFromString(absPath)
		attrs, err := windows.GetFileAttributes(ptr)
		if err != nil {
			return false
		}
		// Check if hidden attribute (0x2) is set
		return attrs&0x2 != 0
	}

	// Unix: check dot prefix
	base := filepath.Base(path)
	return strings.HasPrefix(base, ".")
}
