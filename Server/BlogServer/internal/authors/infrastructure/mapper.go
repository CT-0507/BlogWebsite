package infrastructure

import (
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/authors/domain"
	authordb "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/authors/infrastructure/db"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/model"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/utils"
)

func MapAuthorsAuthorToAuthorProfile(row *authordb.AuthorsAuthor) *domain.AuthorProfile {
	return &domain.AuthorProfile{
		AuthorID:      row.AuthorID,
		UserID:        row.UserID,
		DisplayName:   row.DisplayName,
		Bio:           utils.DerefTextNullable(row.Bio),
		Avatar:        utils.DerefTextNullable(row.Avatar),
		Slug:          row.Slug,
		SocialLink:    utils.DerefTextNullable(row.SocialLink),
		Status:        row.Status,
		Email:         utils.DerefTextNullable(row.Email),
		FollowerCount: row.FollowerCount.Int32,
		BlogCount:     row.BlogCount.Int32,
		Audit: model.Audit{
			CreatedAt: row.CreatedAt.Time,
			CreatedBy: &row.CreatedBy,
			UpdatedAt: row.UpdatedAt.Time,
			UpdatedBy: &row.UpdatedBy,
			DeletedAt: utils.TimePointer(&row.DeletedAt),
			DeletedBy: &row.DeletedBy.String,
		},
	}
}
