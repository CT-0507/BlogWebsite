package mocks

import (
	"context"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/authors/domain"
	"github.com/stretchr/testify/mock"
)

type MockAuthorProfileRepository struct {
	mock.Mock
}

func (m *MockAuthorProfileRepository) CreateAuthorProfile(c context.Context, author *domain.AuthorProfile, userID string, createdBy string) error {
	args := m.Called(c, author, userID, createdBy)

	return args.Error(0)
}

func (m *MockAuthorProfileRepository) GetAuthorProfileByID(
	ctx context.Context,
	authorID string,
	status string,
) (*domain.AuthorProfile, error) {

	args := m.Called(ctx, authorID, status)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.AuthorProfile), args.Error(1)
}

func (m *MockAuthorProfileRepository) GetAuthorProfileBySlug(c context.Context, slug string, status string) (*domain.AuthorProfile, error) {
	args := m.Called(c, slug, status)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.AuthorProfile), args.Error(1)
}

func (m *MockAuthorProfileRepository) ListAuthorProfiles(c context.Context, status string, deletedCheckMode string, page int64, limit int64) (*[]domain.AuthorProfile, error) {
	return nil, nil
}

func (m *MockAuthorProfileRepository) DeleteAuthorProfile(c context.Context, authorID string, userID string) error {

	return nil
}

func (m *MockAuthorProfileRepository) UpdateAuthorStatus(c context.Context, authorID string, status string, userID string) error {

	return nil
}

func (m *MockAuthorProfileRepository) HardDeleteAuthorProfile(c context.Context, authorID string) error {
	return nil
}

func (m *MockAuthorProfileRepository) UpdateAuthorSlug(c context.Context, authorID string, slug string, updatedBy string) error {
	return nil
}

func (m *MockAuthorProfileRepository) CreateAuthorFollower(c context.Context, authorID string, userID string) error {
	args := m.Called(c, authorID, userID)

	return args.Error(0)
}

func (m *MockAuthorProfileRepository) DeleteAuthorFollower(c context.Context, authorID string, userID string) error {
	args := m.Called(c, authorID, userID)

	return args.Error(0)
}

func (m *MockAuthorProfileRepository) GetAuthorFollowers(c context.Context, slug string, page int64, limit int64) ([]string, error) {
	return nil, nil
}

func (m *MockAuthorProfileRepository) GetFollowedAuthors(c context.Context, userID string, page, limit int64) ([]string, error) {
	return nil, nil
}

func (m *MockAuthorProfileRepository) CreateAuthorFeatureBlogs(c context.Context, authorID string, blogIds []string) (int64, error) {
	return 0, nil
}

func (m *MockAuthorProfileRepository) GetAuthorFeaturedBlogIDs(c context.Context, slug string) ([]string, error) {
	return nil, nil
}

func (m *MockAuthorProfileRepository) UpdateAuthorBlogCount(c context.Context, authorID string, isIncrement bool) error {
	return nil
}

func (m *MockAuthorProfileRepository) UpdateAuthorFollowerCount(c context.Context, authorID string, isIncrement bool) error {
	return nil
}
