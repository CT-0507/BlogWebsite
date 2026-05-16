package repository

import (
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
	blogdb "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/infrastructure/db"
)

type BlogRepositoryMapper interface {
	BlogDTOToBlog(blogDTO *blogdb.BlogsBlog) *domain.Blog
	ListBlogsRowDTOToBlog(blogDTO *blogdb.ListBlogsRow) *domain.BlogWithAuthorData
	ListAuthorBlogsByAuthorIDRowDTOToBlog(blogDTO *blogdb.ListBlogsByAuthorRow) *domain.BlogWithAuthorData
	ListAuthorBlogsRowDTOToBlog(blogDTO *blogdb.ListBlogsByAuthorSlugRow) *domain.BlogWithAuthorData
	GetBlogRowDTOToBlogWithAuthorData(blogDTO *blogdb.GetBlogRow) *domain.BlogWithAuthorData
	GetBlogWithReactionDTOToBlogWithAuthorData(blogDTO *blogdb.GetBlogWithUserReactionRow) *domain.BlogWithAuthorData
	GetBlogRowByUrlSlugDTOToBlogWithAuthorData(blogDTO *blogdb.GetBlogByUrlSlugRow) *domain.BlogWithAuthorData
	GetBlogRowDTOToBlog(blogDTO *blogdb.GetBlogRow) *domain.Blog
	MapBlogsCommentToComment(blogComment *blogdb.BlogsComment) *domain.Comment
	MapBlogRootCommentRow(blogComment *blogdb.GetBlogRootCommentRow) (*domain.Comment, error)
	MapBlogRootCommentWithReactionRow(blogComment *blogdb.GetBlogRootCommentWithUserReactionRow) (*domain.Comment, error)
	MapCommentsByParentCommentRow(blogComment *blogdb.GetCommentsByParentCommentRow) (*domain.Comment, error)
	MapCommentsByParentCommentWithReactionRow(blogComment *blogdb.GetCommentsByParentCommentUserWithReactionRow) (*domain.Comment, error)
	MapBlogsIdxUserAuthorProfileToAuthorProfile(author *blogdb.BlogsIdxUserAuthorProfile) *domain.AuthorData
	MapDBBlogReactionToReaction(reaction *blogdb.BlogsBlogReaction) *domain.BlogReaction
	MapDBListRankingRowToRankingBlog(blogDTO *blogdb.ListRankingTableRow) *domain.RankingBlogData
	MapDBReportToBlogReport(report *blogdb.BlogsReport) *domain.BlogReport
}
