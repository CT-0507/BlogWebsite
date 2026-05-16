package application

import (
	"context"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/repository"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
)

type BlogReportUsecases struct {
	txManager database.TxManager
	blogRepo  repository.BlogRepository
}

func NewBlogReportUsecases(
	txManager database.TxManager,
	blogRepo repository.BlogRepository,
) *BlogReportUsecases {
	return &BlogReportUsecases{
		txManager: txManager,
		blogRepo:  blogRepo,
	}
}

func (u *BlogReportUsecases) CreateBlogReport(ctx context.Context, report *domain.BlogReport) (*domain.BlogReport, error) {

	var inserted *domain.BlogReport

	err := u.txManager.WithVoidTx(ctx, func(ctx context.Context) error {

		insertedItem, err := u.blogRepo.InsertBlogReport(ctx, report)
		if err != nil {
			return err
		}

		reportCount, err := u.blogRepo.UpdateBlogReportCount(ctx, inserted.BlogID, 1)
		if err != nil {
			return err
		}

		// Hide blog
		if reportCount > 20 {
			err := u.blogRepo.UpdateBlogStatus(ctx, inserted.BlogID, "over_reported")
			if err != nil {
				return err
			}
		}

		inserted = insertedItem

		return nil
	})

	return inserted, err
}
