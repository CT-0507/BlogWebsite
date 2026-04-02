package infrastructure

import (
	"context"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/authors/domain"
	authordb "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/authors/infrastructure/db"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/utils"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthorProfileRepository struct {
	pool *pgxpool.Pool
}

func NewAuthorProfileRepository(pool *pgxpool.Pool) *AuthorProfileRepository {
	return &AuthorProfileRepository{
		pool: pool,
	}
}

func (r *AuthorProfileRepository) CreateAuthorProfile(c context.Context, author *domain.AuthorProfile, userID string, createdBy string) error {
	db := utils.GetExecutor(c, r.pool)

	q := authordb.New(db)

	return q.CreateAuthorProfile(c, authordb.CreateAuthorProfileParams{
		AuthorID:    author.AuthorID,
		UserID:      userID,
		DisplayName: author.DisplayName,
		Bio: pgtype.Text{
			String: author.DisplayName,
			Valid:  author.DisplayName != "",
		},
		Avatar: pgtype.Text{
			String: author.Avatar,
			Valid:  author.Avatar != "",
		},
		Slug: author.Slug,
		SocialLink: pgtype.Text{
			String: author.SocialLink,
			Valid:  author.SocialLink != "",
		},
		Status: author.Status,
		Email: pgtype.Text{
			String: author.Email,
			Valid:  author.Email != "",
		},
		CreatedBy: createdBy,
	})
}

func (r *AuthorProfileRepository) GetAuthorProfileByID(c context.Context, authorID string, status string, deletedAtCheck string) (*domain.AuthorProfile, error) {
	db := utils.GetExecutor(c, r.pool)

	q := authordb.New(db)

	author, err := q.GetAuthorProfileByID(c, authordb.GetAuthorProfileByIDParams{
		AuthorID: authorID,
		Status:   status,
		Column3:  deletedAtCheck,
	})

	if err != nil {
		return nil, err
	}

	return MapAuthorsAuthorToAuthorProfile(&author), err
}

func (r *AuthorProfileRepository) GetAuthorProfileBySlug(c context.Context, slug string, status string) (*domain.AuthorProfile, error) {
	db := utils.GetExecutor(c, r.pool)

	q := authordb.New(db)

	author, err := q.GetAuthorProfileBySlug(c, authordb.GetAuthorProfileBySlugParams{
		Slug:    slug,
		Status:  status,
		Column3: "check_null",
	})

	if err != nil {
		return nil, err
	}

	return MapAuthorsAuthorToAuthorProfile(&author), err
}

func (r *AuthorProfileRepository) ListAuthorProfiles(c context.Context, status string, deletedCheckMode string, page int64, limit int64) (*[]domain.AuthorProfile, error) {
	db := utils.GetExecutor(c, r.pool)

	q := authordb.New(db)

	rows, err := q.ListAuthorProfies(c, authordb.ListAuthorProfiesParams{
		Status:  "active",
		Column2: "check_null",
	})

	if err != nil {
		return nil, err
	}

	var authorProfiles []domain.AuthorProfile
	for _, value := range rows {
		v := value
		authorProfiles = append(authorProfiles, *MapAuthorsAuthorToAuthorProfile(&v))
	}

	return &authorProfiles, nil
}

func (r *AuthorProfileRepository) DeleteAuthorProfile(c context.Context, authorID string, userID string) error {

	db := utils.GetExecutor(c, r.pool)

	q := authordb.New(db)

	return q.UpdateAuthorProfileDeleteAt(c, authordb.UpdateAuthorProfileDeleteAtParams{
		Status:    "deleted",
		UpdatedBy: userID,
		AuthorID:  authorID,
	})
}

func (r *AuthorProfileRepository) UpdateAuthorStatus(c context.Context, authorID string, status string, userID string) error {

	db := utils.GetExecutor(c, r.pool)

	q := authordb.New(db)

	return q.UpdateAuthorStatus(c, authordb.UpdateAuthorStatusParams{
		Status:    status,
		UpdatedBy: userID,
		AuthorID:  authorID,
	})
}

func (r *AuthorProfileRepository) HardDeleteAuthorProfile(c context.Context, authorID string) error {
	db := utils.GetExecutor(c, r.pool)

	q := authordb.New(db)

	return q.DeleteAuthorProfile(c, authorID)
}

func (r *AuthorProfileRepository) UpdateAuthorSlug(c context.Context, authorID string, slug string, updatedBy string) error {
	db := utils.GetExecutor(c, r.pool)

	q := authordb.New(db)

	return q.UpdateAuthorSlug(c, authordb.UpdateAuthorSlugParams{
		AuthorID:  authorID,
		Slug:      slug,
		UpdatedBy: updatedBy,
	})
}

func (r *AuthorProfileRepository) CreateAuthorFollower(c context.Context, authorID string, userID string) error {
	db := utils.GetExecutor(c, r.pool)

	q := authordb.New(db)

	return q.CreateAuthorFollower(c, authordb.CreateAuthorFollowerParams{
		AuthorID: authorID,
		UserID:   userID,
	})
}

func (r *AuthorProfileRepository) DeleteAuthorFollower(c context.Context, authorID string, userID string) error {
	db := utils.GetExecutor(c, r.pool)

	q := authordb.New(db)

	return q.DeleteAuthorFollower(c, authordb.DeleteAuthorFollowerParams{
		AuthorID: authorID,
		UserID:   userID,
	})
}

func (r *AuthorProfileRepository) GetAuthorFollowers(c context.Context, slug string, page int64, limit int64) ([]string, error) {
	db := utils.GetExecutor(c, r.pool)

	q := authordb.New(db)

	return q.GetAuthorFollowers(c, slug)
}

func (r *AuthorProfileRepository) GetAuthorFollowersByID(c context.Context, authorID string) ([]string, error) {
	db := utils.GetExecutor(c, r.pool)

	q := authordb.New(db)

	return q.GetAuthorFollowersByID(c, authorID)
}

func (r *AuthorProfileRepository) GetFollowedAuthors(c context.Context, userID string, page, limit int64) ([]string, error) {
	db := utils.GetExecutor(c, r.pool)

	q := authordb.New(db)

	return q.GetFollowedAuthors(c, userID)
}

func (r *AuthorProfileRepository) CreateAuthorFeatureBlogs(c context.Context, authorID string, blogIds []string) (int64, error) {
	db := utils.GetExecutor(c, r.pool)

	q := authordb.New(db)

	var params []authordb.CreateAuthorFeatureBlogsParams
	for _, value := range blogIds {
		v := value
		params = append(params, authordb.CreateAuthorFeatureBlogsParams{
			AuthorID: authorID,
			BlogID:   v,
		})
	}

	return q.CreateAuthorFeatureBlogs(c, params)
}

func (r *AuthorProfileRepository) GetAuthorFeaturedBlogIDs(c context.Context, slug string) ([]string, error) {
	db := utils.GetExecutor(c, r.pool)

	q := authordb.New(db)

	return q.GetAuthorFeatureBlogIDs(c, authordb.GetAuthorFeatureBlogIDsParams{
		Slug:   slug,
		Status: "active",
	})
}

func (r *AuthorProfileRepository) UpdateAuthorBlogCount(c context.Context, authorID string, isIncrement bool) error {
	db := utils.GetExecutor(c, r.pool)

	q := authordb.New(db)

	var value int32 = 1
	if !isIncrement {
		value = -1
	}

	return q.UpdateAuthorBlogCount(c, authordb.UpdateAuthorBlogCountParams{
		AuthorID: authorID,
		BlogCount: pgtype.Int4{
			Valid: true,
			Int32: value,
		},
	})
}

func (r *AuthorProfileRepository) UpdateAuthorFollowerCount(c context.Context, authorID string, isIncrement bool) error {
	db := utils.GetExecutor(c, r.pool)

	q := authordb.New(db)

	var value int32 = 1
	if !isIncrement {
		value = -1
	}

	return q.UpdateAuthorFollowerCount(c, authordb.UpdateAuthorFollowerCountParams{
		AuthorID: authorID,
		FollowerCount: pgtype.Int4{
			Valid: true,
			Int32: value,
		},
	})
}
