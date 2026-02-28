package routes

import "github.com/gin-gonic/gin"

type HandlerRoute interface {
	RegisterUnprotectedRoutes(r *gin.Engine)
	RegisterProtectedRoutes(r *gin.Engine)
}
