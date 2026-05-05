package infrastructure

import (
	"context"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
	blogdb "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/infrastructure/db"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/repository"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BlogReactionRepository struct {
	pool   *pgxpool.Pool
	mapper repository.BlogRepositoryMapper
}

func NewBlogReactionRepository(pool *pgxpool.Pool, mapper repository.BlogRepositoryMapper) *BlogReactionRepository {
	return &BlogReactionRepository{
		pool:   pool,
		mapper: mapper,
	}
}

func (r *BlogReactionRepository) UpsertBlogReaction(c context.Context, blogReaction *domain.CreateBlogReaction) (string, string, error) {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	result, err := q.UpsertBlogReaction(c, blogdb.UpsertBlogReactionParams{
		BlogID: blogReaction.BlogID,
		UserID: blogReaction.UserID,
		Type:   blogReaction.Type,
	})

	return result.OldType, result.NewType, err
}
