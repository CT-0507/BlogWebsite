package user

import "github.com/gin-gonic/gin"

func (c *UserHandler) RegisterUnprotectedRoutes(r *gin.Engine) {
	r.POST("/register", c.registerUser)
	r.POST("/login", c.loginUser)
	users := r.Group("/users")
	{
		users.GET("/:user_id", c.getUserById)
	}
}

func (c *UserHandler) RegisterProtectedRoutes(r *gin.Engine) {
	r.POST("/logout", c.logout)
}
