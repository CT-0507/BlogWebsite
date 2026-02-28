package routes

import (
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/middleware"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupProtectedRoutes(r *gin.Engine, pool *pgxpool.Pool, handlers ...HandlerRoute) {
	r.Use(middleware.AuthMiddleWare(pool))

	for _, hanlder := range handlers {
		hanlder.RegisterProtectedRoutes(r)
	}
}
