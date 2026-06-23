package application

import (
	"context"
	"errors"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/repository"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/storage"
)

// List blog related use cases like list all, list with filter, etc
type ListBlogsUseCases struct {
	txManager database.TxManager
	repo      repository.BlogRepository
}

func NewListBlogsUseCases(txManager database.TxManager, repo repository.BlogRepository) *ListBlogsUseCases {
	return &ListBlogsUseCases{
		txManager: txManager,
		repo:      repo,
	}
}

const BLOG_DEFAULT_LIMIT = 10
const (
	MIN_LIMIT = 1
	MAX_LIMIT = 100
)

func (u *ListBlogsUseCases) ListBlogs(ctx context.Context, title, content, author, authorID, sortBy, sortDir *string, page int32, limit int32) (int64, []domain.BlogWithAuthorData, error) {

	offset := (page - 1) * limit

	total, err := u.repo.GetFindAllCount(ctx, title, content, author)
	if err != nil {
		return 0, nil, err
	}
	if total == 0 {
		return total, []domain.BlogWithAuthorData{}, err
	}

	result, err := u.repo.FindAll(ctx, title, content, author, authorID, sortBy, sortDir, offset, limit)
	if err != nil {
		return 0, nil, err
	}

	for i, value := range result {
		if value.ThumbnailUrl != nil {
			thumbnailWithDomain, err := storage.AddDomain(*value.ThumbnailUrl)
			if err != nil {
				return 0, nil, err
			}
			result[i].ThumbnailUrl = &thumbnailWithDomain
		}
	}

	if result != nil && len(result) == 0 {
		return 0, []domain.BlogWithAuthorData{}, nil
	}
	return total, result, nil
}

func (u *ListBlogsUseCases) ListUserLikedBlogs(ctx context.Context, userID string) ([]domain.BlogWithAuthorData, error) {
	result, err := u.repo.GetUserLikedBlogs(ctx, userID)
	if err != nil {
		return nil, err
	}

	for i, value := range result {
		if value.ThumbnailUrl != nil {
			thumbnailWithDomain, err := storage.AddDomain(*value.ThumbnailUrl)
			if err != nil {
				return nil, err
			}
			result[i].ThumbnailUrl = &thumbnailWithDomain
		}
	}

	return result, nil
}

func (u *ListBlogsUseCases) ListBlogsAuthor(ctx context.Context, title, content, sortBy, sortDir *string, page int32, limit int32, userID string) (int64, []domain.BlogWithAuthorData, error) {

	authorID, err := u.repo.VerifyAuthorIDByUserID(ctx, userID)
	if err != nil {
		return 0, nil, err
	}
	if authorID == "" {
		return 0, nil, errors.New("Author not found")
	}

	return u.ListBlogs(ctx, title, content, nil, &authorID, sortBy, sortDir, page, limit)
}

func (u *ListBlogsUseCases) ListAuthorBlogsByAuthorID(ctx context.Context, authorID string) ([]domain.BlogWithAuthorData, error) {
	return u.repo.ListAuthorBlogsByAuthorID(ctx, authorID)
}

func (u *ListBlogsUseCases) ListAuthorBlogsBySlug(ctx context.Context, nickname string) ([]domain.BlogWithAuthorData, error) {
	return u.repo.ListAuthorBlogsBySlug(ctx, nickname)
}

func (u *ListBlogsUseCases) GetRankingBlogsByType(ctx context.Context, searchType string, page, limit int32, shouldGetAll bool, sortBy, sortDir string) (int64, []domain.RankingBlogData, error) {

	offset := (page - 1) * limit
	result, err := u.repo.GetRankingBlogsByType(ctx, searchType, offset, limit, shouldGetAll, sortBy, sortDir)

	if err != nil {
		return 0, nil, err
	}
	if len(result) == 0 {
		return 0, []domain.RankingBlogData{}, nil
	}
	var total int64
	if searchType == "trending" {
		total = *result[0].TotalTrendingResult
	} else {
		total = *result[0].TotalAllTimeResult
	}
	return total, result, nil
}
