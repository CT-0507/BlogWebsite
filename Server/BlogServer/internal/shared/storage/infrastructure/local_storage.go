package infrastructure

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
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

func (l *LocalStorage) Upload(file io.Reader, filename string, contentType string) (string, error) {
	path := filepath.Join(l.BasePath, filename)

	dst, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", l.BaseURL, filename), nil
}

func (l *LocalStorage) MoveFile(src, dst string) error {
	// Open source file
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open src: %w", err)
	}
	defer in.Close()

	// Create destination file
	out, err := os.Create(dst)
	if err != nil {
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
		return fmt.Errorf("copy: %w", err)
	}

	// Sync to disk (optional but safer)
	if err = out.Sync(); err != nil {
		return fmt.Errorf("sync: %w", err)
	}

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
