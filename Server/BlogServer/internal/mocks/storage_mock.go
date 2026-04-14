package mocks

import (
	"io"

	"github.com/stretchr/testify/mock"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) Upload(file io.Reader, filename string, contentType string) (string, error) {

	args := m.Called(file, filename, contentType)

	return "", args.Error(1)
}

func (m *MockStorage) MoveFile(src, dst string) error {
	args := m.Called(src, dst)

	return args.Error(0)
}

func (m *MockStorage) Delete(filename string) error {

	args := m.Called(filename)
	return args.Error(0)
}
