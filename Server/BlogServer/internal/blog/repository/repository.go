package repository

import (
	"context"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
)

type BlogRepository interface {
	Create(c context.Context, blog *domain.Blog) (*domain.Blog, error)
	GetFindAllCount(c context.Context, title, content, author *string) (int64, error)
	FindAll(c context.Context, title, content, author, sortBy, sortDir *string, offset, limit int32) ([]domain.BlogWithAuthorData, error)
	ListAuthorBlogsByAuthorID(c context.Context, authorID string) ([]domain.BlogWithAuthorData, error)
	ListAuthorBlogsBySlug(c context.Context, nickname string) ([]domain.BlogWithAuthorData, error)
	FindByID(c context.Context, id int64) (*domain.BlogWithAuthorData, error)
	FindByUrlSlug(c context.Context, slug string, userID *string) (*domain.BlogWithAuthorData, error)
	// Update(user *Blog) error
	Delete(c context.Context, id int64, userId string) (*int64, error)
	// Cache table
	CreateUserIDAuthorProfileIDCacheRecord(c context.Context, userID string, authorID string, slug string, displayName string) error
	VerifyAuthorIDByUserID(c context.Context, userID string) (string, error)
	// Author deleted event
	UpdateBlogStatusForDeletedAuthor(c context.Context, authorID string) error
	DeleteAuthorHardDeletedBlogs(c context.Context, authorID string) error
	DeleteAuthorCache(c context.Context, authorID string) error
	MarkAuthorCacheAsDeleted(c context.Context, authorID string) error
	RestoreBlog(c context.Context, blogID int64, PreviousStatus string) error
	GetAuthorProfileByUserID(c context.Context, userID string) (*domain.AuthorData, error)
	// Blog metrics
	UpdateBlogReactionCount(c context.Context, blogID int64, transition ReactionTransition) error
	GetRankingBlogsByType(c context.Context, searchType string, offset, limit int32, shouldGetAll bool, sortBy, sortDir string) ([]domain.RankingBlogData, error)
	GetWeeksViews(c context.Context, blogID int64, numberOfWeeks int32) ([]domain.WeekViewData, error)
	GetDaysViews(c context.Context, blogID int64, numberOfDays int32) ([]domain.DateViewData, error)
	UpdateViewCount(c context.Context, blogID int64) error
	// Worker
	UpdateBlogRankingTable(ctx context.Context) error
	TruncateBlogRankingTable(ctx context.Context) error

	// Reports
	UpdateBlogReportCount(c context.Context, blogID int64, delta int64) (int64, error)
	InsertBlogReport(c context.Context, report *domain.BlogReport) (*domain.BlogReport, error)
	DeleteBlogReportByID(c context.Context, reportID int64) (int64, error)
	GetBlogReportsByBlogID(c context.Context, blogID int64) ([]domain.BlogReport, error)
	UpdateBlogStatus(c context.Context, blogID int64, status string) error

	// Author Dashboard
	GetAuthorDashboardViewMetrics(c context.Context, authorID string, userID *string) (*domain.AuthorDashboardViewMetrics, error)
	GetAuthorDashboardReactionMetrics(c context.Context, authorID string, userID *string) (*domain.AuthorDashboardReactionMetrics, error)
}
