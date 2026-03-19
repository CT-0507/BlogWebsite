package infrastructure

import (
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/authors/domain"
	authordb "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/authors/infrastructure/db"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/model"
)

func MapAuthorsAuthorToAuthorProfile(row *authordb.AuthorsAuthor) *domain.AuthorProfile {
	return &domain.AuthorProfile{
		AuthorID:    row.AuthorID,
		UserID:      row.UserID,
		DisplayName: row.DisplayName,
		Bio:         row.Bio.String,
		Avatar:      row.Avatar.String,
		Slug:        row.Slug,
		SocialLink:  row.SocialLink.String,
		Status:      row.Status,
		Email:       row.Email.String,
		Audit: model.Audit{
			CreatedAt: row.CreatedAt.Time,
			CreatedBy: &row.CreatedBy,
			UpdatedAt: row.UpdatedAt.Time,
			UpdatedBy: &row.UpdatedBy,
			DeletedAt: &row.DeletedAt.Time,
			DeletedBy: &row.DeletedBy,
		},
	}
}
