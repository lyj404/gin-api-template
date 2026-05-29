package route

import (
	"github.com/gin-gonic/gin"
	"github.com/lyj404/gin-api-template/api/handler"
	"github.com/lyj404/gin-api-template/api/middleware"
)

func NewRoleRouter(roleHdlr *handler.RoleHandler, rbac *middleware.RBACMiddleware, group *gin.RouterGroup) {
	group.POST("/roles", roleHdlr.CreateRole)
	group.PUT("/roles/:id", roleHdlr.UpdateRole)
	group.DELETE("/roles/:id", roleHdlr.DeleteRole)
	group.GET("/roles/:id", rbac.CheckPermission("role:manage"), roleHdlr.GetRole)
	group.GET("/roles", rbac.CheckPermission("role:manage"), roleHdlr.ListRoles)
	group.POST("/roles/:id/resources", roleHdlr.BindResource)
	group.DELETE("/roles/:id/resources/:resourceId", roleHdlr.UnbindResource)
	group.GET("/roles/:id/resources", rbac.CheckPermission("role:list-resources"), roleHdlr.GetRoleResources)
	group.POST("/roles/:id/menus", roleHdlr.BindMenu)
	group.DELETE("/roles/:id/menus/:menuId", roleHdlr.UnbindMenu)
	group.GET("/roles/:id/menus", roleHdlr.GetRoleMenus)
}
