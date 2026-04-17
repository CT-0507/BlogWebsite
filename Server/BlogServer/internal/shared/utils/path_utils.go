package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

func ToTempPath(src string) (string, error) {
	const prefix = "/uploads/"

	if !strings.HasPrefix(src, prefix) {
		return "", fmt.Errorf("path must start with %s", prefix)
	}

	// already temp
	if strings.HasPrefix(src, "/uploads/temp/") {
		return src, nil
	}

	return strings.Replace(src, prefix, "/uploads/temp/", 1), nil
}

func ToPermanentPath(src string) (string, error) {
	const tempPrefix = "/uploads/temp/"

	if !strings.HasPrefix(src, tempPrefix) {
		return "", fmt.Errorf("path must start with %s", tempPrefix)
	}

	return strings.Replace(src, tempPrefix, "/uploads/", 1), nil
}

func SwapTemp(path string, toTemp bool) string {
	if toTemp {
		return strings.Replace(path, "/uploads/", "/uploads/temp/", 1)
	}
	return strings.Replace(path, "/uploads/temp/", "/uploads/", 1)
}
