package infrastructure

import (
	"context"

	blogdb "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/infrastructure/db"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/repository"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TagRepository struct {
	pool   *pgxpool.Pool
	mapper repository.BlogRepositoryMapper
}

func NewTagRepository(pool *pgxpool.Pool, mapper repository.BlogRepositoryMapper) *TagRepository {
	return &TagRepository{
		pool:   pool,
		mapper: mapper,
	}
}

func (r *TagRepository) UpsertTags(c context.Context, blogID int64, name []string) error {
	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	return q.UpsertTags(c, blogdb.UpsertTagsParams{
		BlogID: blogID,
		Name:   name,
	})
}
