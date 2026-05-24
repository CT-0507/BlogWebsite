package routes

import (
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/middleware"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func SetupProtectedRoutes(r *gin.Engine, pool *pgxpool.Pool, redisClient *redis.Client, handlers ...HandlerRoute) {
	r.Use(middleware.AuthMiddleWare(pool, redisClient))

	for _, hanlder := range handlers {
		hanlder.RegisterProtectedRoutes(r)
	}
}
