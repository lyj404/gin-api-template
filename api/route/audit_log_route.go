package route

import (
	"github.com/gin-gonic/gin"
	"github.com/lyj404/gin-api-template/api/handler"
	"github.com/lyj404/gin-api-template/api/middleware"
)

func NewAuditLogRouter(auditHdlr *handler.AuditLogHandler, rbac *middleware.RBACMiddleware, group *gin.RouterGroup) {
	group.GET("/audit-logs", rbac.CheckPermission("audit:read"), auditHdlr.GetAuditLogsByOperator)
	group.GET("/audit-logs/target", rbac.CheckPermission("audit:read:target"), auditHdlr.GetAuditLogsByTarget)
	group.GET("/audit-logs/time", rbac.CheckPermission("audit:read:time"), auditHdlr.GetAuditLogsByTimeRange)
}
