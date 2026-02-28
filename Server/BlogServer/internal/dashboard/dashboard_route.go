package dashboard

import "github.com/gin-gonic/gin"

func (h *DashboardHandler) RegisterUnprotectedRoutes(r *gin.Engine) {

}

func (h *DashboardHandler) RegisterProtectedRoutes(r *gin.Engine) {
	blogs := r.Group("/dashboard")
	{
		blogs.GET("", h.getDashboardData)
	}
}
