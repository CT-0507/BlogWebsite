package infrastructure

import (
	"context"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
	blogdb "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/infrastructure/db"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/utils"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BlogRepository struct {
	pool *pgxpool.Pool
}

func NewBlogRepository(pool *pgxpool.Pool) *BlogRepository {
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
		UrlSlug:  blog.URLSlug,
		Content: pgtype.Text{
			String: blog.Content,
			Valid:  true,
		},
		CreatedBy: blog.AuthorID,
		UpdatedBy: blog.AuthorID,
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

func (r *BlogRepository) ListAuthorBlogsByAuthorID(c context.Context, authorID string) ([]domain.BlogWithAuthorData, error) {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	rows, err := q.ListBlogsByAuthor(c, blogdb.ListBlogsByAuthorParams{
		AuthorID: authorID,
		Status:   "active",
	})
	if err != nil {
		return nil, err
	}

	var blogs []domain.BlogWithAuthorData
	for _, value := range rows {
		v := value
		blogs = append(blogs, *ListAuthorBlogsByAuthorIDRowDTOToBlog(&v))
	}
	return blogs, nil
}

func (r *BlogRepository) ListAuthorBlogsBySlug(c context.Context, slug string) ([]domain.BlogWithAuthorData, error) {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	rows, err := q.ListBlogsByAuthorSlug(c, blogdb.ListBlogsByAuthorSlugParams{
		Slug:   slug,
		Status: "active",
	})
	if err != nil {
		return nil, err
	}

	var blogs []domain.BlogWithAuthorData
	for _, value := range rows {
		v := value
		blogs = append(blogs, *ListAuthorBlogsRowDTOToBlog(&v))
	}
	return blogs, nil
}

func (r *BlogRepository) FindByID(c context.Context, id int64) (*domain.BlogWithAuthorData, error) {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	row, err := q.GetBlog(c, id)
	if err != nil {
		return nil, err
	}
	return GetBlogRowDTOToBlogWithAuthorData(&row), nil
}

func (r *BlogRepository) FindByUrlSlug(c context.Context, slug string) (*domain.BlogWithAuthorData, error) {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	row, err := q.GetBlogByUrlSlug(c, slug)
	if err != nil {
		return nil, err
	}
	return GetBlogRowByUrlSlugDTOToBlogWithAuthorData(&row), nil
}

// func (r *blogRepository) Update(blog *Blog, q *blogdb.Queries) error {
// 	query := `UPDATE blogs SET name=$1, email=$2 WHERE id=$3`
// 	_, err := r.db.Exec(context.Background(), query, blog.Author, blog.Content, blog.ID)
// 	return err
// }

func (r *BlogRepository) Delete(c context.Context, id int64, userID string) (*int64, error) {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	deletedId, err := q.DeleteBlog(c, blogdb.DeleteBlogParams{
		DeletedBy: pgtype.Text{
			String: userID,
			Valid:  true,
		},
		BlogID: id,
	})
	if err != nil {
		return nil, err
	}
	return &deletedId, nil
}

func (r *BlogRepository) CreateUserIDAuthorProfileIDCacheRecord(c context.Context, userID string, authorID string, slug string, displayName string) error {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	return q.CreateUserAuthorProfileIDCacheRecord(c, blogdb.CreateUserAuthorProfileIDCacheRecordParams{
		UserID:      userID,
		AuthorID:    authorID,
		Slug:        slug,
		DisplayName: displayName,
	})
}

func (r *BlogRepository) VerifyAuthorIDByUserID(c context.Context, userID string) (string, error) {
	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	return q.VerifyAuthorIDByUserID(c, userID)
}

func (r *BlogRepository) UpdateBlogStatusForDeletedAuthor(c context.Context, authorID string) error {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	return q.UpdateBlogStatusForDeletedAuthor(c, authorID)
}

func (r *BlogRepository) DeleteAuthorHardDeletedBlogs(c context.Context, authorID string) error {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	return q.DeleteAuthorHardDeletedBlogs(c, authorID)
}

func (r *BlogRepository) DeleteAuthorCache(c context.Context, authorID string) error {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	return q.DeleteAuthorCache(c, authorID)
}

func (r *BlogRepository) MarkAuthorCacheAsDeleted(c context.Context, authorID string) error {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	return q.MarkAuthorCacheAsDeleted(c, authorID)
}

func (r *BlogRepository) RestoreBlog(c context.Context, blogID int64, PreviousStatus string) error {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	return q.RestoreBlog(c, blogdb.RestoreBlogParams{
		BlogID: blogID,
		Status: PreviousStatus,
	})
}
