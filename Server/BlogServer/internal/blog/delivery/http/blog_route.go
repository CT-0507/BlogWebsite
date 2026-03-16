package http

import "github.com/gin-gonic/gin"

func (h *BlogHandler) RegisterUnprotectedRoutes(r *gin.Engine) {

	blogs := r.Group("/blogs")
	{
		blogs.GET("", h.getAllBlogs)
		blogs.GET("/:slug", h.getBlogByUrlSlug)
		blogs.GET("/id/:id", h.getBlogByID)
	}
}

func (h *BlogHandler) RegisterProtectedRoutes(r *gin.Engine) {

	blogs := r.Group("/blogs")
	{
		blogs.POST("", h.createNewBlog)
		blogs.DELETE("/:id", h.deleteBlogByID)
	}
}
