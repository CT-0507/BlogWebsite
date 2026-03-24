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

	userID, err := utils.GetUserIDStringFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, &gin.H{
			"message": "userId not found",
		})
		return
	}

	if err := h.service.CreateWithOutBox(ctx, &domain.Blog{
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

	blogs, err := h.service.GetAll(ctx)
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

	blogs, err := h.service.ListAuthorBlogsBySlug(ctx, slug)
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
		fmt.Println(parseErr)
		return
	}

	blog, err := h.service.GetBlog(ctx, blogIdInt)
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

	blog, err := h.service.GetBlogByUrlSlug(ctx, slug)
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

	id, err := h.service.DeleteBlog(ctx, blogIdInt, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "blogId not valid",
		})
		return
	}
	c.JSON(http.StatusOK, id)
}
