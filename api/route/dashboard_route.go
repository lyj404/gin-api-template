package route

import (
	"github.com/gin-gonic/gin"
	"github.com/lyj404/gin-api-template/api/handler"
	"github.com/lyj404/gin-api-template/api/middleware"
)

func NewDashboardRouter(h *handler.DashboardHandler, rbac *middleware.RBACMiddleware, group *gin.RouterGroup) {
	group.GET("/dashboard/stats", rbac.CheckPermission("dashboard:read"), h.Stats)
	group.GET("/dashboard/audit-trend", rbac.CheckPermission("dashboard:audit-trend"), h.AuditTrend)
}
