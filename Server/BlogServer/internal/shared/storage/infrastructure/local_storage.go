package infrastructure

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type LocalStorage struct {
	BasePath string
	BaseURL  string
}

func New(
	basePath string,
	baseURL string,
) *LocalStorage {
	return &LocalStorage{
		BasePath: basePath,
		BaseURL:  baseURL,
	}
}

// filename must contain the path it needs to be uploaded
func (l *LocalStorage) Upload(file io.Reader, filename string, contentType string) (string, error) {
	savePath := filepath.Join(l.BasePath, filename)
	log.Println(savePath)
	dst, err := os.Create(savePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", err
	}

	// URL path must use forward slashes
	urlPath := filepath.ToSlash(filename)
	urlPath = strings.TrimPrefix(urlPath, "../")
	return fmt.Sprintf("%s/%s", l.BaseURL, urlPath), nil
}

func (l *LocalStorage) MoveFile(src, dst string) error {

	// remove base URL
	src = l.urlToDiskPath(src)
	dst = l.urlToDiskPath(dst)

	return l.moveFileDisk(src, dst)
}

func (l *LocalStorage) urlToDiskPath(url string) string {

	// remove base URL
	relative := strings.TrimPrefix(url, l.BaseURL)

	// remove leading slash
	relative = strings.TrimPrefix(relative, "/")

	// remove "uploads/" because basePath already points there
	relative = strings.TrimPrefix(relative, "uploads/")

	// build filesystem path
	return filepath.Join(l.BasePath, filepath.FromSlash(relative))
}

func (l *LocalStorage) moveFileDisk(src, dst string) error {
	// Open source file
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open src: %w", err)
	}
	defer in.Close()

	// Create destination file
	out, err := os.Create(dst)
	if err != nil {
		in.Close()
		return fmt.Errorf("create dst: %w", err)
	}
	defer func() {
		// ensure file is closed before checking error
		if cerr := out.Close(); err == nil {
			err = cerr
		}
	}()

	// Copy contents
	if _, err = io.Copy(out, in); err != nil {
		in.Close()
		out.Close()
		return fmt.Errorf("copy: %w", err)
	}

	// Sync to disk (optional but safer)
	if err = out.Sync(); err != nil {
		in.Close()
		out.Close()
		return fmt.Errorf("sync: %w", err)
	}

	in.Close()
	out.Close()

	// Delete source file (like S3 delete after copy)
	if err = os.Remove(src); err != nil {
		return fmt.Errorf("delete src: %w", err)
	}

	return nil
}

func (l *LocalStorage) Delete(filename string) error {
	path := filepath.Join(l.BasePath, filename)
	return os.Remove(path)
}
