package storage

import (
	"context"
	"io"
)

// Local Storage
// type Storage interface {
// 	Upload(file io.Reader, filename string, contentType string) (string, error)
// 	MoveFile(src, dst string) error
// 	Delete(filename string) error
// }

// S3 service interface
type Storage interface {
	// Save file to with temporary flag
	Save(ctx context.Context, key string, body io.Reader, contentType string, isTemporary bool) (*UploadResult, error)
	// Change temporary flag to false
	MarkPermanent(ctx context.Context, key string) error
	MarkDelete(ctx context.Context, key string) error
	// Delete file regardless flag
	Delete(ctx context.Context, key string) error
}
