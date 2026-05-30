package route

import (
	"github.com/gin-gonic/gin"
	"github.com/lyj404/gin-api-template/api/handler"
	"github.com/lyj404/gin-api-template/api/middleware"
)

func NewRoleRouter(roleHdlr *handler.RoleHandler, rbac *middleware.RBACMiddleware, group *gin.RouterGroup) {
	group.POST("/roles", rbac.CheckPermission("role:manage"), roleHdlr.CreateRole)
	group.PUT("/roles/:id", rbac.CheckPermission("role:manage"), roleHdlr.UpdateRole)
	group.DELETE("/roles/:id", rbac.CheckPermission("role:manage"), roleHdlr.DeleteRole)
	group.GET("/roles/:id", rbac.CheckPermission("role:manage"), roleHdlr.GetRole)
	group.GET("/roles", rbac.CheckPermission("role:manage"), roleHdlr.ListRoles)
	group.POST("/roles/:id/resources", rbac.CheckPermission("role:manage"), roleHdlr.BindResource)
	group.DELETE("/roles/:id/resources/:resourceId", rbac.CheckPermission("role:manage"), roleHdlr.UnbindResource)
	group.GET("/roles/:id/resources", rbac.CheckPermission("role:list-resources"), roleHdlr.GetRoleResources)
	group.POST("/roles/:id/menus", rbac.CheckPermission("role:manage"), roleHdlr.BindMenu)
	group.DELETE("/roles/:id/menus/:menuId", rbac.CheckPermission("role:manage"), roleHdlr.UnbindMenu)
	group.GET("/roles/:id/menus", rbac.CheckPermission("role:manage"), roleHdlr.GetRoleMenus)
}
