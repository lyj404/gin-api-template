package route

import (
	"github.com/gin-gonic/gin"
	"github.com/lyj404/gin-api-template/api/handler"
)

func NewDashboardRouter(h *handler.DashboardHandler, group *gin.RouterGroup) {
	group.GET("/dashboard/stats", h.Stats)
	group.GET("/dashboard/audit-trend", h.AuditTrend)
}
