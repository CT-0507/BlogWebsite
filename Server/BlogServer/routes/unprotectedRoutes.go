package routes

import (
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/user"
	"github.com/gin-gonic/gin"
)

func SetupUnprotectedRoutes(r *gin.Engine, blogHandler *blog.BlogHandler, userHandler *user.UserHandler) {
	blogHandler.RegisterUnprotectedRoutes(r)
	userHandler.RegisterUnprotectedRoutes(r)
}
