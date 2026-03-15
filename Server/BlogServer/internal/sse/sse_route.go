package sse

import "github.com/gin-gonic/gin"

func (h *SSEHandler) RegisterUnprotectedRoutes(r *gin.Engine) {
	r.GET("/events/public", h.StreamPublic)
}

func (h *SSEHandler) RegisterProtectedRoutes(r *gin.Engine) {
	r.GET("/events/auth", h.StreamAuth)
}
