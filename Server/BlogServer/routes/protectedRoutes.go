package routes

import (
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/middleware"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/user"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupProtectedRoutes(r *gin.Engine, pool *pgxpool.Pool, blogHandler *blog.BlogHandler, userHandler *user.UserHandler) {
	r.Use(middleware.AuthMiddleWare(pool))

	blogHandler.RegisterProtectedRoutes(r)
	userHandler.RegisterProtectedRoutes(r)
}
