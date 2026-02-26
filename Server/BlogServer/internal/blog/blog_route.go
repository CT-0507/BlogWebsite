package blog

import "github.com/gin-gonic/gin"

func (c *BlogHandler) RegisterUnprotectedRoutes(r *gin.Engine) {
	blogs := r.Group("/blogs")
	{
		blogs.GET("", c.getAllBlogs)
		blogs.GET("/:id", c.getBlogByID)
	}
}

func (c *BlogHandler) RegisterProtectedRoutes(r *gin.Engine) {
	blogs := r.Group("/blogs")
	{
		blogs.POST("", c.createNewBlog)
		blogs.DELETE("/:id", c.deleteBlogByID)
	}
}
