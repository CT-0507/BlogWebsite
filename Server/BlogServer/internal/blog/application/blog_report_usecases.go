package application

import (
	"context"
	"errors"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/repository"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/config"
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

func (u *BlogReportUsecases) GetBlogReportsByBlogID(ctx context.Context, blogID int64, userID string) ([]domain.BlogReport, error) {

	// Confirm identity
	if userID != config.ADMIN_ID && userID != config.SYSTEM_ID {
		authorID, err := u.blogRepo.VerifyAuthorIDByUserID(ctx, userID)
		if err != nil {
			return nil, err
		}
		if authorID == "" {
			return nil, errors.New("Author ID not found")
		}
	}

	return u.blogRepo.GetBlogReportsByBlogID(ctx, blogID)
}
