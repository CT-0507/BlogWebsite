package infrastructure

import (
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
	blogdb "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/infrastructure/db"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/model"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/utils"
	"github.com/jackc/pgx/v5/pgtype"
)

type BlogMapperPG struct{}

func NewBlogMapper() *BlogMapperPG {
	return &BlogMapperPG{}
}

func (b *BlogMapperPG) BlogDTOToBlog(blogDTO *blogdb.BlogsBlog) *domain.Blog {
	return &domain.Blog{
		BlogID:       blogDTO.BlogID,
		Title:        blogDTO.Title,
		URLSlug:      blogDTO.UrlSlug,
		Content:      blogDTO.Content,
		ThumbnailUrl: utils.GetStringPointerFromText(blogDTO.ThumbnailUrl),
		Status:       blogDTO.Status,
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

func (b *BlogMapperPG) ListBlogsRowDTOToBlog(blogDTO *blogdb.ListBlogsRow) *domain.BlogWithAuthorData {
	return &domain.BlogWithAuthorData{
		BlogID:       blogDTO.BlogID,
		Title:        blogDTO.Title,
		URLSlug:      blogDTO.UrlSlug,
		Content:      blogDTO.Content,
		ThumbnailUrl: utils.GetStringPointerFromText(blogDTO.ThumbnailUrl),
		LikeCount:    blogDTO.LikeCount,
		DislikeCount: blogDTO.DislikeCount,
		Status:       blogDTO.Status,
		Tags:         blogDTO.Tags,
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

func (b *BlogMapperPG) ListAuthorBlogsByAuthorIDRowDTOToBlog(blogDTO *blogdb.ListBlogsByAuthorRow) *domain.BlogWithAuthorData {
	return &domain.BlogWithAuthorData{
		BlogID:       blogDTO.BlogID,
		Title:        blogDTO.Title,
		URLSlug:      blogDTO.UrlSlug,
		Content:      blogDTO.Content,
		ThumbnailUrl: utils.GetStringPointerFromText(blogDTO.ThumbnailUrl),
		Status:       blogDTO.Status,
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

func (b *BlogMapperPG) ListAuthorBlogsRowDTOToBlog(blogDTO *blogdb.ListBlogsByAuthorSlugRow) *domain.BlogWithAuthorData {
	return &domain.BlogWithAuthorData{
		BlogID:       blogDTO.BlogID,
		Title:        blogDTO.Title,
		URLSlug:      blogDTO.UrlSlug,
		Content:      blogDTO.Content,
		ThumbnailUrl: utils.GetStringPointerFromText(blogDTO.ThumbnailUrl),
		Status:       blogDTO.Status,
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

func (b *BlogMapperPG) GetBlogRowDTOToBlogWithAuthorData(blogDTO *blogdb.GetBlogRow) *domain.BlogWithAuthorData {
	return &domain.BlogWithAuthorData{
		BlogID:       blogDTO.BlogID,
		Title:        blogDTO.Title,
		URLSlug:      blogDTO.UrlSlug,
		Content:      blogDTO.Content,
		ThumbnailUrl: utils.GetStringPointerFromText(blogDTO.ThumbnailUrl),
		LikeCount:    blogDTO.LikeCount,
		DislikeCount: blogDTO.DislikeCount,
		Status:       blogDTO.Status,
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

func (b *BlogMapperPG) GetBlogWithReactionDTOToBlogWithAuthorData(blogDTO *blogdb.GetBlogWithUserReactionRow) *domain.BlogWithAuthorData {
	return &domain.BlogWithAuthorData{
		BlogID:       blogDTO.BlogID,
		Title:        blogDTO.Title,
		URLSlug:      blogDTO.UrlSlug,
		Content:      blogDTO.Content,
		ThumbnailUrl: utils.GetStringPointerFromText(blogDTO.ThumbnailUrl),
		LikeCount:    blogDTO.LikeCount,
		DislikeCount: blogDTO.DislikeCount,
		UserReaction: utils.GetStringPointerFromText(blogDTO.ReactionType),
		Status:       blogDTO.Status,
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

func (b *BlogMapperPG) GetBlogRowByUrlSlugDTOToBlogWithAuthorData(blogDTO *blogdb.GetBlogByUrlSlugRow) *domain.BlogWithAuthorData {
	return &domain.BlogWithAuthorData{
		BlogID:       blogDTO.BlogID,
		Title:        blogDTO.Title,
		URLSlug:      blogDTO.UrlSlug,
		Content:      blogDTO.Content,
		ThumbnailUrl: utils.GetStringPointerFromText(blogDTO.ThumbnailUrl),
		Status:       blogDTO.Status,
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

func (b *BlogMapperPG) GetBlogRowDTOToBlog(blogDTO *blogdb.GetBlogRow) *domain.Blog {
	return &domain.Blog{
		BlogID:       blogDTO.BlogID,
		Title:        blogDTO.Title,
		Content:      blogDTO.Content,
		ThumbnailUrl: utils.GetStringPointerFromText(blogDTO.ThumbnailUrl),
		Status:       blogDTO.Status,
		Audit: model.Audit{
			CreatedAt: blogDTO.CreatedAt.Time,
			CreatedBy: &blogDTO.CreatedBy,
			UpdatedAt: blogDTO.UpdatedAt.Time,
			UpdatedBy: &blogDTO.UpdatedBy,
		},
	}
}

const DELETED_USER_NAME = "deleted_user"

func (b *BlogMapperPG) getNameOnDeletedActor(id *pgtype.Text) *string {
	if id.Valid == false {
		v := DELETED_USER_NAME
		return &v
	}
	return &id.String
}

func (b *BlogMapperPG) MapBlogsCommentToComment(blogComment *blogdb.BlogsComment) *domain.Comment {
	return &domain.Comment{
		ID:               blogComment.ID,
		BlogID:           blogComment.BlogID,
		ActorType:        blogComment.ActorType,
		ActorID:          b.getNameOnDeletedActor(&blogComment.ActorID),
		ActorDisplayName: blogComment.ActorDisplayName,
		Content:          blogComment.Content,
		Status:           blogComment.Status,
		ParentCommentID:  blogComment.ParentCommentID,
		RootCommentID:    blogComment.RootCommentID,
		Depth:            blogComment.Depth,
		Audit: model.Audit{
			CreatedAt: blogComment.CreatedAt.Time,
			UpdatedAt: blogComment.UpdatedAt.Time,
		},
	}
}

func (b *BlogMapperPG) MapBlogRootCommentRow(blogComment *blogdb.GetBlogRootCommentRow) (*domain.Comment, error) {
	return &domain.Comment{
		ID:               blogComment.ID,
		BlogID:           blogComment.BlogID,
		ActorType:        blogComment.ActorType,
		ActorID:          b.getNameOnDeletedActor(&blogComment.ActorID),
		ActorDisplayName: blogComment.ActorDisplayName,
		Content:          blogComment.Content,
		LikeCount:        blogComment.LikeCount,
		DislikeCount:     blogComment.DislikeCount,
		Status:           blogComment.Status,
		ParentCommentID:  blogComment.ParentCommentID,
		RootCommentID:    blogComment.RootCommentID,
		Depth:            blogComment.Depth,
		ReplyCount:       blogComment.ReplyCount,
		Audit: model.Audit{
			CreatedAt: blogComment.CreatedAt.Time,
			UpdatedAt: blogComment.UpdatedAt.Time,
		},
	}, nil
}

func (b *BlogMapperPG) MapBlogRootCommentWithReactionRow(blogComment *blogdb.GetBlogRootCommentWithUserReactionRow) (*domain.Comment, error) {
	return &domain.Comment{
		ID:               blogComment.ID,
		BlogID:           blogComment.BlogID,
		ActorType:        blogComment.ActorType,
		ActorID:          b.getNameOnDeletedActor(&blogComment.ActorID),
		ActorDisplayName: blogComment.ActorDisplayName,
		Content:          blogComment.Content,
		UserReaction:     utils.GetStringPointerFromText(blogComment.ReactionType),
		LikeCount:        blogComment.LikeCount,
		DislikeCount:     blogComment.DislikeCount,
		Status:           blogComment.Status,
		ParentCommentID:  blogComment.ParentCommentID,
		RootCommentID:    blogComment.RootCommentID,
		Depth:            blogComment.Depth,
		ReplyCount:       blogComment.ReplyCount,
		Audit: model.Audit{
			CreatedAt: blogComment.CreatedAt.Time,
			UpdatedAt: blogComment.UpdatedAt.Time,
		},
	}, nil
}

func (b *BlogMapperPG) MapCommentsByParentCommentRow(blogComment *blogdb.GetCommentsByParentCommentRow) (*domain.Comment, error) {
	return &domain.Comment{
		ID:               blogComment.ID,
		BlogID:           blogComment.BlogID,
		ActorType:        blogComment.ActorType,
		ActorID:          b.getNameOnDeletedActor(&blogComment.ActorID),
		ActorDisplayName: blogComment.ActorDisplayName,
		Content:          blogComment.Content,
		LikeCount:        blogComment.LikeCount,
		DislikeCount:     blogComment.DislikeCount,
		Status:           blogComment.Status,
		ParentCommentID:  blogComment.ParentCommentID,
		RootCommentID:    blogComment.RootCommentID,
		Depth:            blogComment.Depth,
		ReplyCount:       blogComment.ChildCommentCount,
		Audit: model.Audit{
			CreatedAt: blogComment.CreatedAt.Time,
			UpdatedAt: blogComment.UpdatedAt.Time,
		},
	}, nil
}

func (b *BlogMapperPG) MapCommentsByParentCommentWithReactionRow(blogComment *blogdb.GetCommentsByParentCommentUserWithReactionRow) (*domain.Comment, error) {
	return &domain.Comment{
		ID:               blogComment.ID,
		BlogID:           blogComment.BlogID,
		ActorType:        blogComment.ActorType,
		ActorID:          b.getNameOnDeletedActor(&blogComment.ActorID),
		ActorDisplayName: blogComment.ActorDisplayName,
		Content:          blogComment.Content,
		UserReaction:     utils.GetStringPointerFromText(blogComment.ReactionType),
		LikeCount:        blogComment.LikeCount,
		DislikeCount:     blogComment.DislikeCount,
		Status:           blogComment.Status,
		ParentCommentID:  blogComment.ParentCommentID,
		RootCommentID:    blogComment.RootCommentID,
		Depth:            blogComment.Depth,
		ReplyCount:       blogComment.ChildCommentCount,
		Audit: model.Audit{
			CreatedAt: blogComment.CreatedAt.Time,
			UpdatedAt: blogComment.UpdatedAt.Time,
		},
	}, nil
}

func (b *BlogMapperPG) MapBlogsIdxUserAuthorProfileToAuthorProfile(author *blogdb.BlogsIdxUserAuthorProfile) *domain.AuthorData {
	var avatar *string = nil
	if author.Avatar.Valid {
		avatar = &author.Avatar.String
	}
	return &domain.AuthorData{
		AuthorID:    author.AuthorID,
		AvatarURL:   avatar,
		Slug:        author.Slug,
		DisplayName: author.DisplayName,
	}
}

func (b *BlogMapperPG) MapDBBlogReactionToReaction(reaction *blogdb.BlogsBlogReaction) *domain.BlogReaction {
	return &domain.BlogReaction{
		ID:        reaction.ID,
		BlogID:    reaction.BlogID,
		UserID:    reaction.UserID,
		Type:      reaction.Type,
		Status:    reaction.Status,
		CreatedAt: reaction.CreatedAt.Time,
		DeletedAt: utils.TimePointer(&reaction.DeletedAt),
	}
}

func (b *BlogMapperPG) MapDBListRankingRowToRankingBlog(blogDTO *blogdb.ListRankingTableRow) *domain.RankingBlogData {
	return &domain.RankingBlogData{
		BlogID:              blogDTO.BlogID,
		ThumbnailUrl:        utils.GetStringPointerFromText(blogDTO.ThumbnailUrl),
		TotalAllTimeResult:  &blogDTO.TotalAllTime,
		TotalTrendingResult: &blogDTO.TotalTrending,
		RankAllTime:         utils.GetInt32PointerFromInt4(blogDTO.RankAllTime),
		RankTrending:        utils.GetInt32PointerFromInt4(blogDTO.RankTrending),
		ScoreAllTime:        utils.GetFloat64PointerFromFloat8(blogDTO.ScoreAllTime),
		ScoreTrending:       utils.GetFloat64PointerFromFloat8(blogDTO.ScoreTrending),
		LikeCount:           blogDTO.LikeCount,
		DislikeCount:        blogDTO.DislikeCount,
		CommentCount:        blogDTO.CommentCount,
		WeeklyAccessCount:   blogDTO.WeeklyAccessCount,
		DailyAccessCount:    blogDTO.DailyAccessCount,
		CreatedAt:           blogDTO.CreatedAt.Time.String(),
		ComputedAt:          blogDTO.ComputedAt.Time.String(),
		Title:               utils.GetStringPointerFromText(blogDTO.Title),
		AuthorID:            utils.GetStringPointerFromText(blogDTO.AuthorID),
		UrlSlug:             utils.GetStringPointerFromText(blogDTO.UrlSlug),
		Avatar:              utils.GetStringPointerFromText(blogDTO.Avatar),
		DisplayName:         utils.GetStringPointerFromText(blogDTO.DisplayName),
		Slug:                utils.GetStringPointerFromText(blogDTO.Slug),
	}
}
