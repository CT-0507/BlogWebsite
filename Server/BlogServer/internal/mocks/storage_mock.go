package mocks

import (
	"context"
	"io"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/storage"
	"github.com/stretchr/testify/mock"
)

type MockStorage struct {
	mock.Mock
}

// func (m *MockStorage) Upload(file io.Reader, filename string, contentType string) (string, error) {

// 	args := m.Called(file, filename, contentType)

// 	return "", args.Error(1)
// }

// func (m *MockStorage) MoveFile(src, dst string) error {
// 	args := m.Called(src, dst)

// 	return args.Error(0)
// }

// func (m *MockStorage) Delete(filename string) error {

// 	args := m.Called(filename)
// 	return args.Error(0)
// }

func (m *MockStorage) Save(
	ctx context.Context,
	key string,
	body io.Reader,
	contentType string,
	isTemporary bool,
) (*storage.UploadResult, error) {

	args := m.Called(ctx, key, body, contentType, isTemporary)

	return nil, args.Error(1)
}

func (m *MockStorage) MarkPermanent(
	ctx context.Context,
	key string,
) error {
	args := m.Called(ctx, key)

	return args.Error(0)
}

func (m *MockStorage) MarkDelete(
	ctx context.Context,
	key string,
) error {
	args := m.Called(ctx, key)

	return args.Error(0)
}

func (m *MockStorage) Delete(
	ctx context.Context,
	key string,
) error {

	args := m.Called(ctx, key)
	return args.Error(0)
}
