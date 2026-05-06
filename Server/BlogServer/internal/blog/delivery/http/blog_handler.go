package http

import (
	"context"
	"errors"
	"log"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/messages"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateBlogUseCases interface {
	CreateBlogStartSaga(c context.Context, blog *domain.Blog, userID string) error
	VerifyAuthorIDByUserID(c context.Context, userID string) (string, error)
}

type DeleteBlogUseCase interface {
	DeleteBlog(ctx context.Context, id int64, userID string) (*int64, error)
}

type GetBlogUseCases interface {
	GetBlog(ctx context.Context, id int64) (*domain.BlogWithAuthorData, error)
	GetBlogByUrlSlug(ctx context.Context, slug string, userID *string) (*domain.BlogWithAuthorData, error)
}

type ListBlogsUseCases interface {
	ListBlogs(ctx context.Context, title, content, author, sortBy, sortDir *string, page int32, limit int32) (int64, []domain.BlogWithAuthorData, error)
	GetRankingBlogsByType(ctx context.Context, searchType string, page, limit int32, shouldGetAll bool, sortBy, sortDir string) (int64, []domain.RankingBlogData, error)
	ListAuthorBlogsByAuthorID(ctx context.Context, authorID string) ([]domain.BlogWithAuthorData, error)
	ListAuthorBlogsBySlug(ctx context.Context, nickname string) ([]domain.BlogWithAuthorData, error)
}

type CommentUsecases interface {
	CreateComment(c context.Context, newComment *domain.CreateCommentModel, userID string) (*domain.Comment, error)
	GetBlogRootComments(c context.Context, blogID int64, userID *string) (int64, []domain.Comment, error)
	GetChildrenComments(c context.Context, parentCommentID uuid.UUID, userID *string) ([]domain.Comment, error)
	GetCommentByID(c context.Context, commentID uuid.UUID) (*domain.Comment, error)
	HideComment(c context.Context, commentID uuid.UUID, userID string) (int64, error)
	DeleteComment(c context.Context, commentID uuid.UUID, userID string) (int64, error)
	UpdateCommentContent(c context.Context, commentID uuid.UUID, userID string, content string) (int64, error)
}

type CommentReactionUseCases interface {
	CreateCommentReaction(c context.Context, commentReaction *domain.CreateCommentReaction) (int, error)
}
type BlogReactionUseCases interface {
	CreateBlogReaction(c context.Context, blogReaction *domain.CreateBlogReaction) (int, error)
}

type BlogHandler struct {
	createBlogUseCases      CreateBlogUseCases
	getBlogUseCases         GetBlogUseCases
	listBlogsUseCases       ListBlogsUseCases
	deleteBlogUseCases      DeleteBlogUseCase
	commentUsecases         CommentUsecases
	commentReactionUsecases CommentReactionUseCases
	blogReactionUsecases    BlogReactionUseCases
}

func NewBlogHandler(
	createBlogUseCases CreateBlogUseCases,
	getBlogUseCases GetBlogUseCases,
	listBlogsUseCases ListBlogsUseCases,
	deleteBlogUseCases DeleteBlogUseCase,
	commentUsecases CommentUsecases,
	commentReactionUsecases CommentReactionUseCases,
	blogReactionUsecases BlogReactionUseCases,
) *BlogHandler {
	return &BlogHandler{
		createBlogUseCases:      createBlogUseCases,
		getBlogUseCases:         getBlogUseCases,
		listBlogsUseCases:       listBlogsUseCases,
		deleteBlogUseCases:      deleteBlogUseCases,
		commentUsecases:         commentUsecases,
		commentReactionUsecases: commentReactionUsecases,
		blogReactionUsecases:    blogReactionUsecases,
	}
}

// Description: create new blog
//   - @route POST /blogs
//   - @access Private
func (h *BlogHandler) createNewBlog(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c, 2*time.Second)
	defer cancel()

	var blog CreateBlogRequest
	if err := c.ShouldBindJSON(&blog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := utils.ValidateStruct(messages.ENGLISH, blog); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	userID, err := utils.GetUserIDStringFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, &gin.H{
			"message": "userId not found",
		})
		return
	}

	if err := h.createBlogUseCases.CreateBlogStartSaga(ctx, &domain.Blog{
		Title:   blog.Title,
		URLSlug: blog.URLSlug,
		Content: blog.Content,
	}, userID); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, "Okay")
}

// Description: get all blogs
//   - @route GET /blogs
//   - @access Public
func (h *BlogHandler) getAllBlogs(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c, 2*time.Second)
	defer cancel()

	var filter GetBlogFilter
	// ShouldBindQuery binds specifically from query parameters
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if err := utils.ValidateStruct(messages.ENGLISH, filter); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	const (
		DEFAULT_LIMIT = 20
		MIN_LIMIT     = 1
		MAX_LIMIT     = 100
	)

	if filter.Title != nil && *filter.Title == "" {
		filter.Title = nil
	}
	if filter.Content != nil && *filter.Content == "" {
		filter.Content = nil
	}
	if filter.AuthorName != nil && *filter.AuthorName == "" {
		filter.AuthorName = nil
	}

	limitV := int32(DEFAULT_LIMIT)
	if filter.Limit != nil && *filter.Limit >= MIN_LIMIT && *filter.Limit <= MAX_LIMIT {
		limitV = *filter.Limit
	}

	pageV := int32(1)
	if filter.Page != nil && *filter.Page > 0 {
		pageV = *filter.Page
	}

	sortByV := "createdAt"
	if filter.SortBy != nil && *filter.SortBy != "" {
		valid := []string{"title", "relevance", "createdAt"}
		if slices.Contains(valid, *filter.SortBy) {
			sortByV = *filter.SortBy
		}
	}

	sortDirV := "desc"
	if filter.SortDir != nil && *filter.SortDir != "" {
		valid := []string{"asc", "desc"}
		if slices.Contains(valid, *filter.SortDir) {
			sortDirV = *filter.SortDir
		}
	}

	total, blogs, err := h.listBlogsUseCases.ListBlogs(ctx, filter.Title, filter.Content, filter.AuthorName, &sortByV, &sortDirV, pageV, limitV)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"total": total,
		"blogs": &blogs,
	})
}

// Description: get all blogs
//   - @route GET /blogs/author/:slug
//   - @access Public
func (h *BlogHandler) getBlogsByAuthorSlug(c *gin.Context) {
	slug, valid := c.Params.Get("authorSlug")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": messages.MsgRequiredField.FormatLang(messages.ENGLISH, "slug"),
		})
		return
	}

	ctx, cancel := context.WithTimeout(c, 2*time.Second)
	defer cancel()

	blogs, err := h.listBlogsUseCases.ListAuthorBlogsBySlug(ctx, slug)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, blogs)
}

// Description: get blog by id
//   - @route GET /blogs/:id
//   - @access Puclic
func (h *BlogHandler) getBlogByID(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c, 2*time.Second)
	defer cancel()

	blogId, valid := c.Params.Get("id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": messages.MsgRequiredField.FormatLang(messages.ENGLISH, "blogId"),
		})
		return
	}

	blogIdInt, parseErr := strconv.ParseInt(blogId, 10, 64)
	if parseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "blogId not valid",
		})
		log.Println(parseErr)
		return
	}

	blog, err := h.getBlogUseCases.GetBlog(ctx, blogIdInt)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, blog)
}

// Description: get blog by url slug
//   - @route GET /blogs/:slug
//   - @access Puclic
func (h *BlogHandler) getBlogByUrlSlug(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c, 2*time.Second)
	defer cancel()

	slug, valid := c.Params.Get("slug")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": messages.MsgRequiredField.FormatLang(messages.ENGLISH, "slug"),
		})
		return
	}

	var userID *string

	token, err := utils.GetAccessToken(c)
	if token != "" {
		claims, err := utils.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		userID = &claims.UserID
	}

	blog, err := h.getBlogUseCases.GetBlogByUrlSlug(ctx, slug, userID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, blog)
}

// Description: delete blog by id
//   - @route DELETE /blogs/:id
//   - @access Private
func (h *BlogHandler) deleteBlogByID(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c, 2*time.Second)
	defer cancel()

	blogId, valid := c.Params.Get("id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "blogId not found",
		})
		return
	}

	blogIdInt, parseErr := strconv.ParseInt(blogId, 10, 64)
	if parseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "blogId is required",
		})
		return
	}

	userID, error := utils.GetUserIDStringFromContext(c)
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "userID not found",
		})
		return
	}

	id, err := h.deleteBlogUseCases.DeleteBlog(ctx, blogIdInt, userID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, id)
}

func (h *BlogHandler) createComment(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c, 2*time.Second)
	defer cancel()

	blogId, valid := c.Params.Get("id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": messages.MsgRequiredField.FormatLang(messages.ENGLISH, "blogId"),
		})
		return
	}

	blogIdInt, parseErr := strconv.ParseInt(blogId, 10, 64)
	if parseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "blogId not valid",
		})
		return
	}

	var comment CreateCommentRequest
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := utils.ValidateStruct(messages.ENGLISH, comment); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusForbidden, err)
		return
	}

	var parentCommentID *uuid.UUID
	if comment.ParentCommentID != nil {
		v, err := uuid.Parse(*comment.ParentCommentID)
		parentCommentID = &v
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
	}

	var rootCommentID *uuid.UUID
	if comment.RootCommentID != nil {
		v, err := uuid.Parse(*comment.RootCommentID)
		rootCommentID = &v
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
	}

	userAvatar, err := utils.GetAvatarFromContext(c)
	username, err := utils.GetUsernameFromContext(c)
	if err != nil {
		c.AbortWithError(http.StatusForbidden, err)
		return
	}

	insertedComment, err := h.commentUsecases.CreateComment(ctx, &domain.CreateCommentModel{
		BlogID:           blogIdInt,
		ActorType:        comment.ActorType,
		Content:          comment.Content,
		ParentCommentID:  parentCommentID,
		ActorAvatarURL:   userAvatar,
		ActorDisplayName: username,
		RootCommentID:    rootCommentID,
		Depth:            comment.Depth,
	}, userID.String())
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, insertedComment)
}

func (h *BlogHandler) getBlogRootComments(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c, 1*time.Second)
	defer cancel()

	blogId, valid := c.Params.Get("id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": messages.MsgRequiredField.FormatLang(messages.ENGLISH, "blogId"),
		})
		return
	}

	blogIdInt, parseErr := strconv.ParseInt(blogId, 10, 64)
	if parseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "blogId not valid",
		})
		log.Println(parseErr)
		return
	}

	var userID *string

	token, err := utils.GetAccessToken(c)
	if token != "" {
		claims, err := utils.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		userID = &claims.UserID
	}

	total, comments, err := h.commentUsecases.GetBlogRootComments(ctx, blogIdInt, userID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if len(comments) == 0 {
		comments = []domain.Comment{}
	}
	c.JSON(http.StatusOK, gin.H{
		"total":    total,
		"comments": comments,
	})
}

func (h *BlogHandler) getChildrenComments(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c, 1*time.Second)
	defer cancel()

	parentID, valid := c.Params.Get("id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": messages.MsgRequiredField.FormatLang(messages.ENGLISH, "blogId"),
		})
		return
	}

	parentUUID, err := uuid.Parse(parentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "parentID not valid",
		})
		return
	}

	var userID *string

	token, err := utils.GetAccessToken(c)
	if token != "" {
		claims, err := utils.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		userID = &claims.UserID
	}

	comments, err := h.commentUsecases.GetChildrenComments(ctx, parentUUID, userID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if len(comments) == 0 {
		comments = []domain.Comment{}
	}
	c.JSON(http.StatusOK, comments)
}

func (h *BlogHandler) getCommentByID(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c, 1*time.Second)
	defer cancel()

	id, valid := c.Params.Get("id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": messages.MsgRequiredField.FormatLang(messages.ENGLISH, "blogId"),
		})
		return
	}

	uuid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "commentID is not valid",
		})
		return
	}

	comment, err := h.commentUsecases.GetCommentByID(ctx, uuid)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, comment)
}

func (h *BlogHandler) HideCommentByID(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c, 1*time.Second)
	defer cancel()

	id, valid := c.Params.Get("id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": messages.MsgRequiredField.FormatLang(messages.ENGLISH, "blogId"),
		})
		return
	}

	uuid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "commentID is not valid",
		})
		return
	}

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusForbidden, err)
		return
	}

	_, err = h.commentUsecases.HideComment(ctx, uuid, userID.String())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"commentId": id,
	})
}

func (h *BlogHandler) DeleteCommentByID(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c, 1*time.Second)
	defer cancel()

	id, valid := c.Params.Get("id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": messages.MsgRequiredField.FormatLang(messages.ENGLISH, "blogId"),
		})
		return
	}

	uuid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "commentID is not valid",
		})
		return
	}

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusForbidden, err)
		return
	}

	_, err = h.commentUsecases.DeleteComment(ctx, uuid, userID.String())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"commentId": id,
	})
}

func (h *BlogHandler) UpdateCommentContentByID(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c, 1*time.Second)
	defer cancel()

	id, valid := c.Params.Get("id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": messages.MsgRequiredField.FormatLang(messages.ENGLISH, "blogId"),
		})
		return
	}

	uuid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "commentID is not valid",
		})
		return
	}

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusForbidden, err)
		return
	}

	var content UpdateCommentContentRequest
	if err := c.ShouldBindJSON(&content); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := utils.ValidateStruct(messages.ENGLISH, content); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	_, err = h.commentUsecases.UpdateCommentContent(ctx, uuid, userID.String(), content.Content)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"commentId": id,
		"content":   content.Content,
	})
}

func (h *BlogHandler) CreateBlogReaction(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c, 1*time.Second)
	defer cancel()

	blogId, valid := c.Params.Get("id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": messages.MsgRequiredField.FormatLang(messages.ENGLISH, "blogId"),
		})
		return
	}

	blogIdInt, parseErr := strconv.ParseInt(blogId, 10, 64)
	if parseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "blogId not valid",
		})
		log.Println(parseErr)
		return
	}

	var reaction CreateBlogReactionRequest
	if err := c.ShouldBindJSON(&reaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := utils.ValidateStruct(messages.ENGLISH, reaction); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	reactionMap := []string{"like", "dislike"}
	if !slices.Contains(reactionMap, reaction.Type) {
		c.JSON(http.StatusBadRequest, errors.New("Invalid reaction type."))
		return
	}

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusForbidden, err)
		return
	}

	transitionType, err := h.blogReactionUsecases.CreateBlogReaction(ctx, &domain.CreateBlogReaction{
		Type:   reaction.Type,
		BlogID: blogIdInt,
		UserID: userID.String(),
	})

	transtionMap := map[int]string{
		0: "AddLike",
		1: "AddDislike",
		2: "LikeToDislike",
		3: "DislikeToLike",
	}

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"transitionType": transtionMap[transitionType],
		"blogId":         blogIdInt,
		"type":           reaction.Type,
	})
}

func (h *BlogHandler) CreateCommentReaction(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c, 1*time.Second)
	defer cancel()

	commentId, valid := c.Params.Get("id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": messages.MsgRequiredField.FormatLang(messages.ENGLISH, "commentId"),
		})
		return
	}

	commentUUID, parseErr := uuid.Parse(commentId)
	if parseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "blogId not valid",
		})
		log.Println(parseErr)
		return
	}

	var reaction CreateCommentReactionRequest
	if err := c.ShouldBindJSON(&reaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := utils.ValidateStruct(messages.ENGLISH, reaction); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	reactionMap := []string{"like", "dislike"}
	if !slices.Contains(reactionMap, reaction.Type) {
		c.JSON(http.StatusBadRequest, errors.New("Invalid reaction type."))
		return
	}

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusForbidden, err)
		return
	}

	transitionType, err := h.commentReactionUsecases.CreateCommentReaction(ctx, &domain.CreateCommentReaction{
		Type:      reaction.Type,
		CommentID: commentUUID,
		UserID:    userID.String(),
	})

	transtionMap := map[int]string{
		0: "AddLike",
		1: "AddDislike",
		2: "LikeToDislike",
		3: "DislikeToLike",
	}

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"transitionType": transtionMap[transitionType],
		"commentId":      commentId,
		"type":           reaction.Type,
	})
}

func (h *BlogHandler) GetRankingBlogsByType(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c, 2*time.Second)
	defer cancel()

	var filter GetBlogRankingFilter
	// ShouldBindQuery binds specifically from query parameters
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if err := utils.ValidateStruct(messages.ENGLISH, filter); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	const (
		DEFAULT_LIMIT = 20
		MIN_LIMIT     = 1
		MAX_LIMIT     = 100
	)

	searchType := "createdAt"
	if filter.Type != nil && *filter.Type != "" {
		valid := []string{"allTime", "trending"}
		if slices.Contains(valid, *filter.Type) {
			searchType = *filter.Type
		}
	}

	limitV := int32(DEFAULT_LIMIT)
	if filter.Limit != nil && *filter.Limit >= MIN_LIMIT && *filter.Limit <= MAX_LIMIT {
		limitV = *filter.Limit
	}

	pageV := int32(1)
	if filter.Page != nil && *filter.Page > 0 {
		pageV = *filter.Page
	}

	sortByV := "createdAt"
	if filter.SortBy != nil && *filter.SortBy != "" {
		valid := []string{"title", "relevance", "createdAt"}
		if slices.Contains(valid, *filter.SortBy) {
			sortByV = *filter.SortBy
		}
	}

	sortDirV := "desc"
	if filter.SortDir != nil && *filter.SortDir != "" {
		valid := []string{"asc", "desc"}
		if slices.Contains(valid, *filter.SortDir) {
			sortDirV = *filter.SortDir
		}
	}

	total, blogs, err := h.listBlogsUseCases.GetRankingBlogsByType(ctx, searchType, pageV, limitV, false, sortByV, sortDirV)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"total": total,
		"blogs": &blogs,
	})
}
