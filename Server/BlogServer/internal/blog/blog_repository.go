package blog

import (
	"context"

	blogdb "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type BlogRepository interface {
	Create(c context.Context, q *blogdb.Queries, blog *Blog) (*Blog, error)
	FindAll(c context.Context, q *blogdb.Queries) ([]BlogWithAuthorDTO, error)
	FindByID(c context.Context, q *blogdb.Queries, id int64) (*Blog, error)
	// Update(user *Blog, q *blogdb.Queries) error
	Delete(c context.Context, q *blogdb.Queries, id int64, userId uuid.UUID) (*int64, error)
}

type blogRepository struct {
}

func NewBlogRepository() BlogRepository {
	return &blogRepository{}
}

func (r *blogRepository) Create(c context.Context, q *blogdb.Queries, blog *Blog) (*Blog, error) {

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

func (r *blogRepository) FindAll(c context.Context, q *blogdb.Queries) ([]BlogWithAuthorDTO, error) {
	rows, err := q.ListBlogs(c)
	if err != nil {
		return nil, err
	}

	var blogs []BlogWithAuthorDTO
	for _, value := range rows {
		v := value
		blogs = append(blogs, *ListBlogsRowDTOToBlog(&v))
	}
	return blogs, nil
}

func (r *blogRepository) FindByID(c context.Context, q *blogdb.Queries, id int64) (*Blog, error) {
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

func (r *blogRepository) Delete(c context.Context, q *blogdb.Queries, id int64, userId uuid.UUID) (*int64, error) {
	deletedId, err := q.DeleteBlog(c, blogdb.DeleteBlogParams{
		DeletedBy: &userId,
		BlogID:    id,
	})
	if err != nil {
		return nil, err
	}
	return &deletedId, nil
}
