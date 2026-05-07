package worker

import (
	"context"
	"time"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
)

type BlogRankingOperator interface {
	UpdateBlogRankingTable(ctx context.Context) error
	TruncateBlogRankingTable(ctx context.Context) error
}

type BlogWorker struct {
	txManager           database.TxManager
	blogRankingOperator BlogRankingOperator
}

func NewBlogWorker(txManager database.TxManager, blogRankingOperator BlogRankingOperator) *BlogWorker {
	return &BlogWorker{
		txManager:           txManager,
		blogRankingOperator: blogRankingOperator,
	}
}

func (w *BlogWorker) StartUpdateRankingTable(ctx context.Context) {

	ticker := time.NewTicker(2 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.updateBlogRankingTable(ctx)

		case <-ctx.Done():
			return
		}
	}
}

func (w *BlogWorker) updateBlogRankingTable(ctx context.Context) error {
	return w.txManager.WithVoidTx(ctx, func(ctx context.Context) error {

		err := w.blogRankingOperator.TruncateBlogRankingTable(ctx)
		if err != nil {
			return err
		}

		return w.blogRankingOperator.UpdateBlogRankingTable(ctx)
	})
}
