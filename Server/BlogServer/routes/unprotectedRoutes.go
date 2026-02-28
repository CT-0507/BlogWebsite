package routes

import (
	"github.com/gin-gonic/gin"
)

func SetupUnprotectedRoutes(r *gin.Engine, handlers ...HandlerRoute) {
	for _, handler := range handlers {
		handler.RegisterUnprotectedRoutes(r)
	}
}
