package application

import (
	"context"
	"errors"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/repository"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
)

type BlogMetricsUsecases struct {
	txManager database.TxManager
	blogRepo  repository.BlogRepository
}

func NewBlogMetricsUsecases(
	txManager database.TxManager,
	blogRepo repository.BlogRepository,
) *BlogMetricsUsecases {
	return &BlogMetricsUsecases{
		txManager: txManager,
		blogRepo:  blogRepo,
	}
}

func (u *BlogMetricsUsecases) GetWeeksViews(ctx context.Context, blogID int64, numberOfDays int32) ([]domain.WeekViewData, error) {
	return u.blogRepo.GetWeeksViews(ctx, blogID, numberOfDays)
}

func (u *BlogMetricsUsecases) GetDateViews(ctx context.Context, blogID int64, numberOfDays int32) ([]domain.DateViewData, error) {
	return u.blogRepo.GetDaysViews(ctx, blogID, numberOfDays)
}

func (u *BlogMetricsUsecases) GetAuthorDashboardMetrics(ctx context.Context, authorID *string, userID *string) (*domain.AuthorDashboardViewMetrics, *domain.AuthorDashboardReactionMetrics, error) {

	if authorID == nil && userID == nil {
		return nil, nil, errors.New("authorID or userID is required")
	}

	var authorIDV *string = authorID
	if userID != nil {
		authorIDQuery, err := u.blogRepo.VerifyAuthorIDByUserID(ctx, *userID)
		if err != nil {
			return nil, nil, err
		}
		if authorIDQuery == "" {
			return nil, nil, errors.New("Author not found")
		}
		authorIDV = &authorIDQuery
	}

	viewMetrics, err := u.blogRepo.GetAuthorDashboardViewMetrics(ctx, *authorIDV, userID)
	if err != nil {
		return nil, nil, err
	}

	reactionMetrics, err := u.blogRepo.GetAuthorDashboardReactionMetrics(ctx, *authorIDV, userID)
	if err != nil {
		return nil, nil, err
	}

	return viewMetrics, reactionMetrics, nil
}
