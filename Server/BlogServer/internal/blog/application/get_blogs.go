package application

import (
	"context"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/repository"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
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

func (s *ListBlogsUseCases) ListBlogs(ctx context.Context, title, content, author, sortBy, sortDir *string, page int32, limit int32) (int64, []domain.BlogWithAuthorData, error) {

	offset := (page - 1) * limit

	total, err := s.repo.GetFindAllCount(ctx, title, content, author)
	if err != nil {
		return 0, nil, err
	}
	if total == 0 {
		return total, []domain.BlogWithAuthorData{}, err
	}

	result, err := s.repo.FindAll(ctx, title, content, author, sortBy, sortDir, offset, limit)
	if err != nil {
		return 0, nil, err
	}
	if result != nil && len(result) == 0 {
		return 0, []domain.BlogWithAuthorData{}, nil
	}
	return total, result, nil
}

func (s *ListBlogsUseCases) ListAuthorBlogsByAuthorID(ctx context.Context, authorID string) ([]domain.BlogWithAuthorData, error) {
	return s.repo.ListAuthorBlogsByAuthorID(ctx, authorID)
}

func (s *ListBlogsUseCases) ListAuthorBlogsBySlug(ctx context.Context, nickname string) ([]domain.BlogWithAuthorData, error) {
	return s.repo.ListAuthorBlogsBySlug(ctx, nickname)
}

func (s *ListBlogsUseCases) GetRankingBlogsByType(ctx context.Context, searchType string, page, limit int32, shouldGetAll bool, sortBy, sortDir string) (int64, []domain.RankingBlogData, error) {

	offset := (page - 1) * limit
	result, err := s.repo.GetRankingBlogsByType(ctx, searchType, offset, limit, shouldGetAll, sortBy, sortDir)

	if err != nil {
		return 0, nil, err
	}
	if result != nil && len(result) == 0 {
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
