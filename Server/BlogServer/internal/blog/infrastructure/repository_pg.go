package infrastructure

import (
	"context"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
	blogdb "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/infrastructure/db"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/repository"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/utils"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BlogRepository struct {
	pool   *pgxpool.Pool
	mapper repository.BlogRepositoryMapper
}

func NewBlogRepository(pool *pgxpool.Pool, mapper repository.BlogRepositoryMapper) *BlogRepository {
	return &BlogRepository{
		pool:   pool,
		mapper: mapper,
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

	return r.mapper.BlogDTOToBlog(&newBlog), nil
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
		blogs = append(blogs, *r.mapper.ListBlogsRowDTOToBlog(&v))
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
		blogs = append(blogs, *r.mapper.ListAuthorBlogsByAuthorIDRowDTOToBlog(&v))
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
		blogs = append(blogs, *r.mapper.ListAuthorBlogsRowDTOToBlog(&v))
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
	return r.mapper.GetBlogRowDTOToBlogWithAuthorData(&row), nil
}

func (r *BlogRepository) FindByUrlSlug(c context.Context, slug string, userID *string) (*domain.BlogWithAuthorData, error) {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	if userID != nil {
		row, err := q.GetBlogWithUserReaction(c, blogdb.GetBlogWithUserReactionParams{
			UserID:  *userID,
			UrlSlug: slug,
		})
		if err != nil {
			return nil, err
		}
		return r.mapper.GetBlogWithReactionDTOToBlogWithAuthorData(&row), nil
	} else {
		row, err := q.GetBlogByUrlSlug(c, slug)
		if err != nil {
			return nil, err
		}
		return r.mapper.GetBlogRowByUrlSlugDTOToBlogWithAuthorData(&row), nil
	}
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

func (r *BlogRepository) GetAuthorProfileByUserID(c context.Context, userID string) (*domain.AuthorData, error) {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	author, err := q.GetAuthorCacheByUserID(c, userID)
	if err != nil {
		return nil, err
	}

	return r.mapper.MapBlogsIdxUserAuthorProfileToAuthorProfile(&author), nil
}

func (r *BlogRepository) UpdateBlogReactionCount(c context.Context, blogID int64, transition repository.ReactionTransition) error {

	db := utils.GetExecutor(c, r.pool)

	q := blogdb.New(db)

	var likeDelta int32 = 0
	var dislikeDelta int32 = 0

	switch transition {
	case repository.AddLike:
		likeDelta++
	case repository.AddDislike:
		dislikeDelta++
	case repository.LikeToDislike:
		likeDelta--
		dislikeDelta++
	case repository.DislikeToLike:
		likeDelta++
		dislikeDelta--
	}

	return q.UpdateBlogReactionCount(c, blogdb.UpdateBlogReactionCountParams{
		LikeCount:    likeDelta,
		DislikeCount: dislikeDelta,
		BlogID:       blogID,
	})
}
