package user

import "github.com/gin-gonic/gin"

func (h *UserHandler) RegisterUnprotectedRoutes(r *gin.Engine) {
	r.POST("/register", h.registerUser)
	r.POST("/login", h.loginUser)
	r.POST("/refresh", h.RefreshTokenHandler)
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
		user.POST("/change-email-code", h.GetChangeEmailCode)
		user.POST("/change-email", h.UpdateUserEmail)
		user.POST("/change-basic-info", h.UpdateUserBasicInfo)
		user.POST("/change-password", h.ChangePassword)
	}
}
