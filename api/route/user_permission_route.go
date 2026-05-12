package route

import (
	"github.com/gin-gonic/gin"
	"github.com/lyj404/gin-api-template/api/handler"
)

// NewUserPermissionRouter 注册用户权限相关路由
func NewUserPermissionRouter(userPermHdlr *handler.UserPermissionHandler, group *gin.RouterGroup) {
	user := group.Group("/user")
	{
		user.GET("/permissions", userPermHdlr.GetUserPermissions)
		user.GET("/menus", userPermHdlr.GetUserMenus)
	}
}
