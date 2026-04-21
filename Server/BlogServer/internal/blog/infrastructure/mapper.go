package infrastructure

import (
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
	blogdb "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/infrastructure/db"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/model"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/utils"
	"github.com/jackc/pgx/v5/pgtype"
)

func BlogDTOToBlog(blogDTO *blogdb.BlogsBlog) *domain.Blog {
	return &domain.Blog{
		BlogID:  blogDTO.BlogID,
		Title:   blogDTO.Title,
		URLSlug: blogDTO.UrlSlug,
		Content: blogDTO.Content.String,
		Status:  blogDTO.Status,
		Audit: model.Audit{
			CreatedAt: blogDTO.CreatedAt.Time,
			CreatedBy: &blogDTO.CreatedBy,
			UpdatedAt: blogDTO.UpdatedAt.Time,
			UpdatedBy: &blogDTO.UpdatedBy,
			DeletedAt: utils.TimePointer(&blogDTO.DeletedAt),
			DeletedBy: &blogDTO.DeletedBy.String,
		},
	}
}

func ListBlogsRowDTOToBlog(blogDTO *blogdb.ListBlogsRow) *domain.BlogWithAuthorData {
	return &domain.BlogWithAuthorData{
		BlogID:  blogDTO.BlogID,
		Title:   blogDTO.Title,
		URLSlug: blogDTO.UrlSlug,
		Content: blogDTO.Content.String,
		Status:  blogDTO.Status,
		Author: domain.AuthorData{
			AuthorID:    blogDTO.AuthorID,
			Slug:        blogDTO.Slug,
			DisplayName: blogDTO.DisplayName,
		},
		Audit: model.Audit{
			CreatedAt: blogDTO.CreatedAt.Time,
			CreatedBy: &blogDTO.CreatedBy,
			UpdatedAt: blogDTO.UpdatedAt.Time,
			UpdatedBy: &blogDTO.UpdatedBy,
		},
	}
}

func ListAuthorBlogsByAuthorIDRowDTOToBlog(blogDTO *blogdb.ListBlogsByAuthorRow) *domain.BlogWithAuthorData {
	return &domain.BlogWithAuthorData{
		BlogID:  blogDTO.BlogID,
		Title:   blogDTO.Title,
		URLSlug: blogDTO.UrlSlug,
		Content: blogDTO.Content.String,
		Status:  blogDTO.Status,
		Author: domain.AuthorData{
			AuthorID:    blogDTO.AuthorID,
			Slug:        blogDTO.Slug,
			DisplayName: blogDTO.DisplayName,
		},
		Audit: model.Audit{
			CreatedAt: blogDTO.CreatedAt.Time,
			CreatedBy: &blogDTO.CreatedBy,
			UpdatedAt: blogDTO.UpdatedAt.Time,
			UpdatedBy: &blogDTO.UpdatedBy,
		},
	}
}

func ListAuthorBlogsRowDTOToBlog(blogDTO *blogdb.ListBlogsByAuthorSlugRow) *domain.BlogWithAuthorData {
	return &domain.BlogWithAuthorData{
		BlogID:  blogDTO.BlogID,
		Title:   blogDTO.Title,
		URLSlug: blogDTO.UrlSlug,
		Content: blogDTO.Content.String,
		Status:  blogDTO.Status,
		Author: domain.AuthorData{
			AuthorID:    blogDTO.AuthorID,
			Slug:        blogDTO.Slug,
			DisplayName: blogDTO.DisplayName,
		},
		Audit: model.Audit{
			CreatedAt: blogDTO.CreatedAt.Time,
			CreatedBy: &blogDTO.CreatedBy,
			UpdatedAt: blogDTO.UpdatedAt.Time,
			UpdatedBy: &blogDTO.UpdatedBy,
		},
	}
}

func GetBlogRowDTOToBlogWithAuthorData(blogDTO *blogdb.GetBlogRow) *domain.BlogWithAuthorData {
	return &domain.BlogWithAuthorData{
		BlogID:  blogDTO.BlogID,
		Title:   blogDTO.Title,
		URLSlug: blogDTO.UrlSlug,
		Content: blogDTO.Content.String,
		Status:  blogDTO.Status,
		Author: domain.AuthorData{
			AuthorID:    blogDTO.AuthorID,
			Slug:        blogDTO.Slug,
			DisplayName: blogDTO.DisplayName,
		},
		Audit: model.Audit{
			CreatedAt: blogDTO.CreatedAt.Time,
			CreatedBy: &blogDTO.CreatedBy,
			UpdatedAt: blogDTO.UpdatedAt.Time,
			UpdatedBy: &blogDTO.UpdatedBy,
		},
	}
}

func GetBlogRowByUrlSlugDTOToBlogWithAuthorData(blogDTO *blogdb.GetBlogByUrlSlugRow) *domain.BlogWithAuthorData {
	return &domain.BlogWithAuthorData{
		BlogID:  blogDTO.BlogID,
		Title:   blogDTO.Title,
		URLSlug: blogDTO.UrlSlug,
		Content: blogDTO.Content.String,
		Status:  blogDTO.Status,
		Author: domain.AuthorData{
			AuthorID:    blogDTO.AuthorID,
			Slug:        blogDTO.Slug,
			DisplayName: blogDTO.DisplayName,
		},
		Audit: model.Audit{
			CreatedAt: blogDTO.CreatedAt.Time,
			CreatedBy: &blogDTO.CreatedBy,
			UpdatedAt: blogDTO.UpdatedAt.Time,
			UpdatedBy: &blogDTO.UpdatedBy,
		},
	}
}

func GetBlogRowDTOToBlog(blogDTO *blogdb.GetBlogRow) *domain.Blog {
	return &domain.Blog{
		BlogID:  blogDTO.BlogID,
		Title:   blogDTO.Title,
		Content: blogDTO.Content.String,
		Status:  blogDTO.Status,
		Audit: model.Audit{
			CreatedAt: blogDTO.CreatedAt.Time,
			CreatedBy: &blogDTO.CreatedBy,
			UpdatedAt: blogDTO.UpdatedAt.Time,
			UpdatedBy: &blogDTO.UpdatedBy,
		},
	}
}

const DELETED_USER_NAME = "deleted_user"

func getNameOnDeletedActor(id *pgtype.Text) *string {
	if id.Valid == false {
		v := DELETED_USER_NAME
		return &v
	}
	return &id.String
}

func MapBlogsCommentToComment(blogComment *blogdb.BlogsComment) *domain.Comment {
	return &domain.Comment{
		ID:               blogComment.ID,
		BlogID:           blogComment.BlogID,
		ActorType:        blogComment.ActorType,
		ActorID:          getNameOnDeletedActor(&blogComment.ActorID),
		ActorDisplayName: blogComment.ActorDisplayName,
		Content:          blogComment.Content,
		Status:           blogComment.Status,
		ParentCommentID:  blogComment.ParentCommentID,
		RootCommentID:    blogComment.RootCommentID,
		Depth:            blogComment.Depth,
	}
}
