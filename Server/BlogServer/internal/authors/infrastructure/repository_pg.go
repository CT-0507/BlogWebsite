package infrastructure

import (
	"context"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/authors/domain"
	authordb "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/authors/infrastructure/db"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthorProfileRepository struct {
	pool *pgxpool.Pool
}

func NewAuthorProfileRepository(pool *pgxpool.Pool) domain.AuthorProfileRepository {
	return &AuthorProfileRepository{
		pool: pool,
	}
}

func (r *AuthorProfileRepository) CreateAuthorProfile(c context.Context, author *domain.AuthorProfile, userID uuid.UUID, createdBy uuid.UUID) error {
	db := utils.GetExecutor(c, r.pool)

	q := authordb.New(db)

	return q.CreateAuthorProfile(c, authordb.CreateAuthorProfileParams{
		AuthorID:    author.AuthorID,
		UserID:      &userID,
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
		UpdatedBy: &createdBy,
	})
}

func (r *AuthorProfileRepository) ListAuthorProfies(c context.Context, status string, deletedCheckMode string) ([]domain.AuthorProfile, error) {
	db := utils.GetExecutor(c, r.pool)

	q := authordb.New(db)

	rows, err := q.ListAuthorProfies(c, authordb.ListAuthorProfiesParams{
		Status:  "active",
		Column2: "check_not_null",
	})

	if err != nil {
		return nil, err
	}

	var authorProfiles []domain.AuthorProfile
	for _, value := range rows {
		v := value
		authorProfiles = append(authorProfiles, *MapAuthorsAuthorToAuthorProfileList(&v))
	}

	return authorProfiles, nil
}
