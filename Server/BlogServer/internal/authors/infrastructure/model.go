package infrastructure

import (
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/authors/domain"
	authordb "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/authors/infrastructure/db"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/model"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/utils"
)

func MapAuthorsAuthorToAuthorProfileList(row *authordb.AuthorsAuthor) *domain.AuthorProfile {
	return &domain.AuthorProfile{
		AuthorID:    row.AuthorID,
		UserID:      row.UserID.String(),
		DisplayName: row.DisplayName,
		Bio:         row.Bio.String,
		Avatar:      row.Avatar.String,
		Slug:        row.Slug,
		SocialLink:  row.SocialLink.String,
		Status:      row.Status,
		Email:       row.Email.String,
		Audit: model.Audit{
			CreatedAt: row.CreatedAt.Time,
			CreatedBy: utils.UUIDPtr(row.CreatedBy),
			UpdatedAt: row.UpdatedAt.Time,
			UpdatedBy: utils.UUIDPtr(row.UpdatedBy),
			DeletedAt: utils.TimePointer(&row.DeletedAt),
			DeletedBy: utils.UUIDPtr(row.DeletedBy),
		},
	}
}
