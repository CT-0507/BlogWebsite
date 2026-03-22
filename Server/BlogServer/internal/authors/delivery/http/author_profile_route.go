package http

import "github.com/gin-gonic/gin"

func (h *AuthorProfileHandler) RegisterUnprotectedRoutes(r *gin.Engine) {

	api := r.Group("/api")
	v1 := api.Group("/v1/authors")
	{
		v1.GET("", h.listAuthorProfiles)
		v1.GET("/:slug", h.getAuthorProfileBySlug)
		v1.GET("/id/:id", h.getAuthorProfileByID)
		v1.GET("/:slug/followers", h.getAuthorFollowers)
	}
}

func (h *AuthorProfileHandler) RegisterProtectedRoutes(r *gin.Engine) {

	api := r.Group("/api")
	v1 := api.Group("/v1/authors")
	{
		// identity
		v1.POST("", h.createAuthorProfile)
		v1.DELETE("/:id", h.deleteAuthorProfile)
		v1.DELETE("/:id/hard", h.hardDeleteAuthorProfile)
		v1.PATCH("/:id/slug", h.updateAuthorSlug)
		v1.PATCH("/:id/status", h.updateAuthorStatus)

		// follower system
		v1.POST("/:id/follow", h.followAuthor)
		v1.DELETE("/:id/follow", h.unfollowAuthor)
		v1.GET("me/following/authors", h.getFollowedAuthors)

		// Presentation
		v1.PUT("/:slug/featured-blogs", h.setAuthorFeaturedBlogs)
		v1.GET("/:slug/featured-blogs", h.getAuthorFeaturedBlogs)
	}
}
