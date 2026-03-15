package blog

import (
	blogdb "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog_LayeredArchitecture/db"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/model"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/utils"
)

func BlogDTOToBlog(blogDTO *blogdb.BlogsBlog) *Blog {
	return &Blog{
		BlogID:  blogDTO.BlogID,
		Title:   blogDTO.Title,
		Content: blogDTO.Content.String,
		Active:  blogDTO.Active,
		Audit: model.Audit{
			CreatedAt: blogDTO.CreatedAt.Time,
			CreatedBy: utils.UUIDPtr(blogDTO.CreatedBy),
			UpdatedAt: blogDTO.UpdatedAt.Time,
			UpdatedBy: utils.UUIDPtr(blogDTO.UpdatedBy),
			DeletedAt: utils.TimePointer(&blogDTO.DeletedAt),
			DeletedBy: utils.UUIDPtr(blogDTO.DeletedBy),
		},
	}
}

func ListBlogsRowDTOToBlog(blogDTO *blogdb.ListBlogsRow) *BlogWithAuthorDTO {
	return &BlogWithAuthorDTO{
		BlogID:  blogDTO.BlogID,
		Title:   blogDTO.Title,
		Content: blogDTO.Content.String,
		Active:  blogDTO.Active,
		Author: AuthorDTO{
			AuthorID: blogDTO.AuthorID.String(),
			FullName: blogDTO.AuthorName.(string),
		},
		Audit: model.Audit{
			CreatedAt: blogDTO.CreatedAt.Time,
			CreatedBy: utils.UUIDPtr(blogDTO.CreatedBy),
			UpdatedAt: blogDTO.UpdatedAt.Time,
			UpdatedBy: utils.UUIDPtr(blogDTO.UpdatedBy),
		},
	}
}

func GetBlogRowDTOToBlog(blogDTO *blogdb.GetBlogRow) *Blog {
	return &Blog{
		BlogID:  blogDTO.BlogID,
		Title:   blogDTO.Title,
		Content: blogDTO.Content.String,
		Active:  blogDTO.Active,
		Audit: model.Audit{
			CreatedAt: blogDTO.CreatedAt.Time,
			CreatedBy: utils.UUIDPtr(blogDTO.CreatedBy),
			UpdatedAt: blogDTO.UpdatedAt.Time,
			UpdatedBy: utils.UUIDPtr(blogDTO.UpdatedBy),
		},
	}
}
