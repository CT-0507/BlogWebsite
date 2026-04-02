package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// EnsureUploadPath creates (if not exists) and returns a Y/M/D folder path
func EnsureUploadPath(baseDir string) (string, error) {
	now := time.Now()

	// build path: baseDir/YYYY/MM/DD
	path := filepath.Join(
		baseDir,
		fmt.Sprintf("%d", now.Year()),
		fmt.Sprintf("%02d", int(now.Month())),
		fmt.Sprintf("%02d", now.Day()),
	)

	// create directory if not exists
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return "", err
	}

	return path, nil
}
