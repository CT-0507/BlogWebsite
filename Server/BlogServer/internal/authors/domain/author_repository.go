package domain

import (
	"context"
)

type AuthorProfileRepository interface {

	// Author identity
	CreateAuthorProfile(c context.Context, author *AuthorProfile, userID string, createdBy string) error
	GetAuthorProfileByID(c context.Context, authorID string, status string) (*AuthorProfile, error)
	GetAuthorProfileBySlug(c context.Context, slug string, status string) (*AuthorProfile, error)
	ListAuthorProfiles(c context.Context, status string, deletedCheckMode string, page int64, limit int64) (*[]AuthorProfile, error)
	// UpdateAuthorProfile(c context.Context, author *AuthorProfile, userID string) error

	// Author Slug & Identity Management
	// Soft delete
	DeleteAuthorProfile(c context.Context, authorID string, userID string) error
	// Hard delete
	HardDeleteAuthorProfile(c context.Context, authorID string) error
	UpdateAuthorStatus(c context.Context, authorID string, status string, userID string) error
	UpdateAuthorSlug(c context.Context, authorID string, slug string, updatedBy string) error

	// Follower system
	CreateAuthorFollower(c context.Context, authorID string, userID string) error
	DeleteAuthorFollower(c context.Context, authorID string, userID string) error
	GetAuthorFollowers(c context.Context, slug string, page int64, limit int64) ([]string, error)
	GetFollowedAuthors(c context.Context, userID string, page int64, limit int64) ([]string, error)

	// Author Social (Follow System)
	CreateAuthorFeatureBlogs(c context.Context, authorID string, blogIds []string) (int64, error) // Max 10
	GetAuthorFeaturedBlogIDs(c context.Context, slug string) ([]string, error)
	// UpdateAuthorFeatureBlog()
	// Stats system

	// Author discovery
	// SearchAuthors(c context.Context, keyword string, pagination int) ([]*AuthorProfile, error)
	// GetTrendingAuthors
	UpdateAuthorBlogCount(c context.Context, authorID string, isIncrement bool) error
	UpdateAuthorFollowerCount(c context.Context, authorID string, isIncrement bool) error
}
