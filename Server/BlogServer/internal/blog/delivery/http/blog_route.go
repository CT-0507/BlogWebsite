package http

import "github.com/gin-gonic/gin"

func (h *BlogHandler) RegisterUnprotectedRoutes(r *gin.Engine) {

	api := r.Group("/api")
	v1 := api.Group("/v1")

	blogs := v1.Group("/blogs")
	blogs.GET("", h.getAllBlogs)
	blogs.GET("/ranking", h.getRankingBlogsByType)
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

}

func (h *BlogHandler) RegisterProtectedRoutes(r *gin.Engine) {

	api := r.Group("/api")
	v1 := api.Group("/v1")

	authorDashboard := v1.Group("/dashboard/author")
	authorDashboard.GET("", h.getAuthorDashboardMetrics)
	blogsDashboard := authorDashboard.Group("/blogs")
	blogsDashboard.GET("", h.getAllBlogsAuthor)
	blogDashboard := blogsDashboard.Group("/:id")
	blogDashboard.GET("/metrics", h.getViewsData)

	blogs := v1.Group("/blogs")
	blogs.POST("", h.createNewBlog)
	blogs.GET("/me/liked-blogs", h.getUserLikedBlogs)
	blog := blogs.Group("/:id")
	blog.DELETE("", h.deleteBlogByID)
	blog.POST("/reaction", h.CreateBlogReaction)
	blog.PATCH("", h.updateNewBlog)

	reports := blog.Group("/reports")
	reports.GET("", h.getBlogReportsByBlogID)
	reports.POST("", h.createBlogReport)

	comments := blog.Group("/comments")
	comments.POST("", h.createComment)

	sComments := v1.Group("/comments")
	sComment := sComments.Group("/:id")
	sComment.POST("/reaction", h.createCommentReaction)
	sComment.DELETE("/delete", h.deleteCommentByID)
	sComment.PATCH("/hidden", h.hideCommentByID)
	sComment.PATCH("", h.updateCommentContentByID)

	r.MaxMultipartMemory = 8 << 20
	uploads := v1.Group("/uploads")
	uploads.POST("/image", h.uploadImage)
}
