package storage

import "io"

// Shared
type FileStorageParams struct {
	File        io.Reader
	FileName    string
	ContentType string
}
