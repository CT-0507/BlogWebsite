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

func (l *LocalStorage) Delete(filename string) error {
	path := filepath.Join(l.BasePath, filename)
	return os.Remove(path)
}
