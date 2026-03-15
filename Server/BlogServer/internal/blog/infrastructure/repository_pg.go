package infrastructure

import (
	"context"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
	blogdb "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/infrastructure/db"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BlogRepository struct {
	pool *pgxpool.Pool
}

func NewBlogRepository(pool *pgxpool.Pool) domain.BlogRepository {
	return &BlogRepository{
		pool: pool,
	}
}

func (r *BlogRepository) Create(c context.Context, blog *domain.Blog) (*domain.Blog, error) {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	newBlog, err := q.CreateBlog(c, blogdb.CreateBlogParams{
		AuthorID: blog.AuthorID,
		Title:    blog.Title,
		Content: pgtype.Text{
			String: blog.Content,
			Valid:  true,
		},
		CreatedBy: &blog.AuthorID,
		UpdatedBy: &blog.AuthorID,
	})

	if err != nil {
		return nil, err
	}

	return BlogDTOToBlog(&newBlog), nil
}

func (r *BlogRepository) FindAll(c context.Context) ([]domain.BlogWithAuthorData, error) {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	rows, err := q.ListBlogs(c)
	if err != nil {
		return nil, err
	}

	var blogs []domain.BlogWithAuthorData
	for _, value := range rows {
		v := value
		blogs = append(blogs, *ListBlogsRowDTOToBlog(&v))
	}
	return blogs, nil
}

func (r *BlogRepository) FindByID(c context.Context, id int64) (*domain.Blog, error) {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	row, err := q.GetBlog(c, id)
	if err != nil {
		return nil, err
	}
	return GetBlogRowDTOToBlog(&row), nil
}

// func (r *blogRepository) Update(blog *Blog, q *blogdb.Queries) error {
// 	query := `UPDATE blogs SET name=$1, email=$2 WHERE id=$3`
// 	_, err := r.db.Exec(context.Background(), query, blog.Author, blog.Content, blog.ID)
// 	return err
// }

func (r *BlogRepository) Delete(c context.Context, id int64, userId uuid.UUID) (*int64, error) {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	deletedId, err := q.DeleteBlog(c, blogdb.DeleteBlogParams{
		DeletedBy: &userId,
		BlogID:    id,
	})
	if err != nil {
		return nil, err
	}
	return &deletedId, nil
}
