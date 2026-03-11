package sse

import "github.com/gin-gonic/gin"

func (h *SSEHandler) RegisterUnprotectedRoutes(r *gin.Engine) {

}

func (h *SSEHandler) RegisterProtectedRoutes(r *gin.Engine) {
	r.GET("/events", h.Stream)
}
