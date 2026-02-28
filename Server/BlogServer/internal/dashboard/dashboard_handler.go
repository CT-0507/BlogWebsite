package dashboard

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
}

func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{}
}

func (h *DashboardHandler) getDashboardData(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to dashboard"})
}
