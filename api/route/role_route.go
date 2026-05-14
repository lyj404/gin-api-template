package route

import (
	"github.com/gin-gonic/gin"
	"github.com/lyj404/gin-api-template/api/handler"
)

func NewRoleRouter(roleHdlr *handler.RoleHandler, group *gin.RouterGroup) {
	group.POST("/roles", roleHdlr.CreateRole)
	group.PUT("/roles/:id", roleHdlr.UpdateRole)
	group.DELETE("/roles/:id", roleHdlr.DeleteRole)
	group.GET("/roles/:id", roleHdlr.GetRole)
	group.GET("/roles", roleHdlr.ListRoles)
	group.POST("/roles/:id/resources", roleHdlr.BindResource)
	group.DELETE("/roles/:id/resources/:resourceId", roleHdlr.UnbindResource)
	group.GET("/roles/:id/resources", roleHdlr.GetRoleResources)
}
