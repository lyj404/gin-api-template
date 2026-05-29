package route

import (
	"github.com/gin-gonic/gin"
	"github.com/lyj404/gin-api-template/api/handler"
	"github.com/lyj404/gin-api-template/api/middleware"
)

func NewResourceRouter(resourceHdlr *handler.ResourceHandler, rbac *middleware.RBACMiddleware, group *gin.RouterGroup) {
	group.POST("/resources", resourceHdlr.CreateResource)
	group.PUT("/resources/:id", resourceHdlr.UpdateResource)
	group.DELETE("/resources/:id", resourceHdlr.DeleteResource)
	group.GET("/resources/:id", rbac.CheckPermission("resource:manage"), resourceHdlr.GetResource)
	group.GET("/resources", rbac.CheckPermission("resource:manage"), resourceHdlr.ListResources)
}
