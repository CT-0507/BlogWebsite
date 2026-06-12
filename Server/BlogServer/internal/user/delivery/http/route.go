package http

import "github.com/gin-gonic/gin"

func (h *UserHandler) RegisterUnprotectedRoutes(r *gin.Engine) {

	api := r.Group("/api")
	v1 := api.Group("/v1")

	v1.POST("/register", h.registerUser)
	v1.POST("/login", h.loginUser)
	v1.POST("/refresh", h.refreshToken)
	v1.GET("/get-hashed-string", h.getHashedString)
	users := r.Group("/users")
	users.GET("/:user_id", h.getUserById)

	v1.POST("/contact/new", h.createContactForm)
}

func (h *UserHandler) RegisterProtectedRoutes(r *gin.Engine) {

	api := r.Group("/api")
	v1 := api.Group("/v1")

	v1.DELETE("/contact/:id", h.deleteContactForm)

	v1.GET("/me", h.getUserById)
	v1.POST("/logout", h.logout)

	user := v1.Group("/user")
	user.GET("/notifications", h.getNotifications)
	user.POST("/change-email-code", h.getChangeEmailCode)
	user.POST("/change-email", h.updateUserEmail)
	user.POST("/change-basic-info", h.updateUserBasicInfo)
	user.POST("/change-password", h.changePassword)
}
