package storage

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Generate a path with folder name and yyyy/mm/dd format.
//
// Ex: GenerateKey("rte", "cat.png")
//
// Result: rte/2026/06/14/cat.png
func GenerateKey(folder, originalFilename string) string {
	ext := strings.ToLower(filepath.Ext(originalFilename))

	now := time.Now().UTC()

	return fmt.Sprintf(
		"%s/%04d/%02d/%02d/%s%s",
		folder,
		now.Year(),
		now.Month(),
		now.Day(),
		uuid.NewString(),
		ext,
	)
}

func StripURL(inputURL string) string {
	// Handle URLs without scheme
	if !strings.HasPrefix(inputURL, "http://") &&
		!strings.HasPrefix(inputURL, "https://") {
		inputURL = "https://" + inputURL
	}

	u, err := url.Parse(inputURL)
	if err != nil {
		return ""
	}

	return u.Path
}

func AddDomain(inputURL string) (string, error) {

	domain := os.Getenv("CLOUD_FRONT_DOMAIN")
	if domain == "" {
		return "", errors.New("Domain not found")
	}

	return fmt.Sprintf("https://%s%s", domain, inputURL), nil
}
