package application

import (
	"context"

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
