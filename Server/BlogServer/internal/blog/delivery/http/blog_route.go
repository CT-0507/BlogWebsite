package http

import "github.com/gin-gonic/gin"

func (h *BlogHandler) RegisterUnprotectedRoutes(r *gin.Engine) {

	api := r.Group("/api")
	v1 := api.Group("/v1")

	blogs := v1.Group("/blogs")
	blogs.GET("", h.getAllBlogs)
	blogs.GET("/ranking", h.GetRankingBlogsByType)
	blogs.GET("/slug/:slug", h.getBlogByUrlSlug)
	blogs.GET("/author/slug/:authorSlug", h.getBlogsByAuthorSlug)
	blog := blogs.Group("/:id")
	blog.GET("", h.getBlogByID)

	blogComments := blog.Group("/comments")
	blogComments.GET("", h.getBlogRootComments)

	comments := v1.Group("/comments")
	comment := comments.Group("/:id")
	comment.GET("/children", h.getChildrenComments)
	comment.GET("", h.getCommentByID)

	blog.GET("/metrics", h.GetViewsData)
}

func (h *BlogHandler) RegisterProtectedRoutes(r *gin.Engine) {

	api := r.Group("/api")
	v1 := api.Group("/v1")

	blogs := v1.Group("/blogs")
	blogs.POST("", h.createNewBlog)
	blog := blogs.Group("/:id")
	blog.DELETE("", h.deleteBlogByID)
	blog.POST("/reaction", h.CreateBlogReaction)

	comments := blog.Group("/comments")
	comments.POST("", h.createComment)

	sComments := v1.Group("/comments")
	sComment := sComments.Group("/:id")
	sComment.POST("/reaction", h.CreateCommentReaction)
	sComment.DELETE("/delete", h.DeleteCommentByID)
	sComment.PATCH("/hidden", h.HideCommentByID)
	sComment.PATCH("", h.UpdateCommentContentByID)
}
