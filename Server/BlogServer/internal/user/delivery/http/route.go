package http

import "github.com/gin-gonic/gin"

func (h *UserHandler) RegisterUnprotectedRoutes(r *gin.Engine) {
	r.POST("/register", h.RegisterUser)
	r.POST("/login", h.LoginUser)
	r.POST("/refresh", h.RefreshTokenHandler)
	r.GET("/get-hashed-string", h.GetHashedString)
	users := r.Group("/users")
	{
		users.GET("/:user_id", h.getUserById)
	}
}

func (h *UserHandler) RegisterProtectedRoutes(r *gin.Engine) {
	r.GET("/me", h.getUserById)
	r.POST("/logout", h.logout)
	user := r.Group("/user")
	{
		user.GET("/notifications", h.GetNotifications)
		user.POST("/change-email-code", h.GetChangeEmailCode)
		user.POST("/change-email", h.UpdateUserEmail)
		user.POST("/change-basic-info", h.UpdateUserBasicInfo)
		user.POST("/change-password", h.ChangePassword)
	}
}
