package infrastructure

import (
	"context"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
	blogdb "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/infrastructure/db"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/repository"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CommentReactionRepository struct {
	pool   *pgxpool.Pool
	mapper repository.BlogRepositoryMapper
}

func NewCommentReactionRepository(pool *pgxpool.Pool, mapper repository.BlogRepositoryMapper) *CommentReactionRepository {
	return &CommentReactionRepository{
		pool:   pool,
		mapper: mapper,
	}
}

func (r *CommentReactionRepository) UpsertCommentReaction(c context.Context, blogReaction *domain.CreateCommentReaction) (string, string, error) {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	result, err := q.UpsertCommentReaction(c, blogdb.UpsertCommentReactionParams{
		CommentID: blogReaction.CommentID,
		UserID:    blogReaction.UserID,
		Type:      blogReaction.Type,
	})

	return result.OldType, result.NewType, err
}
