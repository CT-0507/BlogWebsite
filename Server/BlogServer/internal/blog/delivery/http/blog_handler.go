package http

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/application"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/messages"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BlogHandler struct {
	service application.BlogService
}

func NewBlogHandler(service application.BlogService) *BlogHandler {
	return &BlogHandler{service: service}
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

	uuid, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, &gin.H{
			"message": "userId not found",
		})
		return
	}

	if err := h.service.CreateWithOutBox(ctx, &domain.Blog{
		AuthorID: uuid,
		Title:    blog.Title,
		Content:  blog.Content,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, "Okay")
}

func (h *BlogHandler) getAllBlogs(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	blogs, err := h.service.GetAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(blogs) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "success with no data",
		})
		return
	}
	c.JSON(http.StatusOK, blogs)
}

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
		fmt.Println(parseErr)
		return
	}

	blog, err := h.service.GetBlog(ctx, blogIdInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "blogId not valid",
		})
		return
	}
	c.JSON(http.StatusOK, blog)
}

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

	userId, valid := c.Params.Get("user_id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "blogId not found",
		})
		return
	}

	userUUId, pErr := uuid.Parse(userId)
	if pErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "userId is invalid",
		})
		return
	}

	id, err := h.service.DeleteBlog(ctx, blogIdInt, userUUId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "blogId not valid",
		})
		return
	}
	c.JSON(http.StatusOK, id)
}
