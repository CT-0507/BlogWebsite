package storage

import "io"

// Provides needed parameters for saving
type FileStorageParams struct {
	// File content
	File io.Reader
	// File name. Ex: cat.png
	FileName string
	// Ex: .png
	ContentType string
}

type UploadResult struct {
	Key string `json:"key"`
	URL string `json:"url"`
}
