package http

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/messages"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/utils"
	"github.com/gin-gonic/gin"
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
	GetBlogByUrlSlug(ctx context.Context, slug string) (*domain.BlogWithAuthorData, error)
}

type ListBlogsUseCases interface {
	ListBlogs(ctx context.Context) ([]domain.BlogWithAuthorData, error)
	ListAuthorBlogsByAuthorID(ctx context.Context, authorID string) ([]domain.BlogWithAuthorData, error)
	ListAuthorBlogsBySlug(ctx context.Context, nickname string) ([]domain.BlogWithAuthorData, error)
}

type BlogHandler struct {
	createBlogUseCases CreateBlogUseCases
	getBlogUseCases    GetBlogUseCases
	listBlogsUseCases  ListBlogsUseCases
	deleteBlogUseCases DeleteBlogUseCase
}

func NewBlogHandler(
	createBlogUseCases CreateBlogUseCases,
	getBlogUseCases GetBlogUseCases,
	listBlogsUseCases ListBlogsUseCases,
	deleteBlogUseCases DeleteBlogUseCase,
) *BlogHandler {
	return &BlogHandler{
		createBlogUseCases: createBlogUseCases,
		getBlogUseCases:    getBlogUseCases,
		listBlogsUseCases:  listBlogsUseCases,
		deleteBlogUseCases: deleteBlogUseCases,
	}
}

// Description: create new blog
//   - @route POST /blogs
//   - @access Private
func (h *BlogHandler) createNewBlog(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c, 10*time.Second)
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, "Okay")
}

// Description: get all blogs
//   - @route GET /blogs
//   - @access Public
func (h *BlogHandler) getAllBlogs(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	blogs, err := h.listBlogsUseCases.ListBlogs(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, blogs)
}

// Description: get all blogs
//   - @route GET /blogs/author/:slug
//   - @access Public
func (h *BlogHandler) getBlogsByAuthorSlug(c *gin.Context) {
	slug, valid := c.Params.Get("slug")
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, blogs)
}

// Description: get blog by id
//   - @route GET /blogs/id/:id
//   - @access Puclic
func (h *BlogHandler) getBlogByID(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c, 10*time.Second)
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "blogId not found",
		})
		return
	}
	c.JSON(http.StatusOK, blog)
}

// Description: get blog by url slug
//   - @route GET /blogs/:slug
//   - @access Puclic
func (h *BlogHandler) getBlogByUrlSlug(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	slug, valid := c.Params.Get("slug")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": messages.MsgRequiredField.FormatLang(messages.ENGLISH, "slug"),
		})
		return
	}

	blog, err := h.getBlogUseCases.GetBlogByUrlSlug(ctx, slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Url not found",
		})
		return
	}
	c.JSON(http.StatusOK, blog)
}

// Description: delete blog by id
//   - @route DELETE /blogs/:id
//   - @access Private
func (h *BlogHandler) deleteBlogByID(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c, 10*time.Second)
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "blogId not valid",
		})
		return
	}
	c.JSON(http.StatusOK, id)
}
