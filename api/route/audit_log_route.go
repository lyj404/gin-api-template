package route

import (
	"github.com/gin-gonic/gin"
	"github.com/lyj404/gin-api-template/api/handler"
)

func NewAuditLogRouter(auditHdlr *handler.AuditLogHandler, group *gin.RouterGroup) {
	group.GET("/audit-logs", auditHdlr.GetAuditLogsByOperator)
	group.GET("/audit-logs/target", auditHdlr.GetAuditLogsByTarget)
	group.GET("/audit-logs/time", auditHdlr.GetAuditLogsByTimeRange)
}
