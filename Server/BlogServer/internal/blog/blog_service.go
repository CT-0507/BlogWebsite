package blog

import (
	"context"
	"time"

	blogdb "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BlogService interface {
	Create(c context.Context, blog *Blog) error
	GetAll(c context.Context) ([]BlogWithAuthorDTO, error)
	GetByID(c context.Context, id int64) (*Blog, error)
	// Update(blog *Blog) error
	Delete(c context.Context, id int64, userId uuid.UUID) (*int64, error)
}

type blogService struct {
	pool *pgxpool.Pool
	repo BlogRepository
}

func NewBlogService(pool *pgxpool.Pool, repo BlogRepository) BlogService {
	return &blogService{
		pool: pool,
		repo: repo,
	}
}

func (s *blogService) withTx(
	ctx context.Context,
	fn func(q *blogdb.Queries) error,
) error {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	queries := blogdb.New(tx)

	if err := fn(queries); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

//	func (s *blogService) Create(c context.Context, blog *Blog) error {
//		// if blog.Author == "" {
//		// 	return errors.New("name is required")
//		// }
//		ctx, cancel := context.WithTimeout(c, 10*time.Second)
//		defer cancel()
//		q := blogdb.New(s.pool)
//		_, err := s.repo.Create(ctx, q, blog)
//		return err
//	}
func (s *blogService) Create(c context.Context, blog *Blog) error {
	return s.withTx(c, func(q *blogdb.Queries) error {
		_, err := s.repo.Create(c, q, blog)
		return err
	})
}

func (s *blogService) GetAll(c context.Context) ([]BlogWithAuthorDTO, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()
	q := blogdb.New(s.pool)
	return s.repo.FindAll(ctx, q)
}

func (s *blogService) GetByID(c context.Context, id int64) (*Blog, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()
	q := blogdb.New(s.pool)
	return s.repo.FindByID(ctx, q, id)
}

// func (s *blogService) Update(user *Blog) error {
// 	return s.repo.Update(user)
// }

func (s *blogService) Delete(c context.Context, id int64, userId uuid.UUID) (*int64, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()
	q := blogdb.New(s.pool)
	return s.repo.Delete(ctx, q, id, userId)
}
