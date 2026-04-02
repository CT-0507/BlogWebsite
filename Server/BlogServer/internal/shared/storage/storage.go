package storage

import "io"

type Storage interface {
	Upload(file io.Reader, filename string, contentType string) (string, error)
	Delete(filename string) error
}
