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
	v1 := api.Group("/v1")
	authors := v1.Group("/authors")
	// identity
	authors.POST("", h.createAuthorProfile)
	authors.DELETE("/:id", h.deleteAuthorProfile)
	authors.DELETE("/:id/hard", h.hardDeleteAuthorProfile)
	authors.PATCH("/:id/slug", h.updateAuthorSlug)
	authors.PATCH("/:id/status", h.updateAuthorStatus)

	// follower system
	authors.POST("/:id/follow", h.followAuthor)
	authors.DELETE("/:id/follow", h.unfollowAuthor)

	// Presentation
	authors.PUT("/:slug/featured-blogs", h.setAuthorFeaturedBlogs)
	authors.GET("/:slug/featured-blogs", h.getAuthorFeaturedBlogs)

	// Me
	me := v1.Group("/me")
	authors.GET("/me/following/authors", h.getFollowedAuthors)
	me.GET("/authorProfile", h.getMyAuthorProfile)
}
