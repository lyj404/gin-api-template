package route

import (
	"github.com/gin-gonic/gin"
	"github.com/lyj404/gin-api-template/api/handler"
	"github.com/lyj404/gin-api-template/api/middleware"
)

func NewMenuRouter(menuHandler *handler.MenuHandler, rbac *middleware.RBACMiddleware, group *gin.RouterGroup) {
	menus := group.Group("/menus")
	{
		menus.POST("", menuHandler.CreateMenu)
		menus.GET("", rbac.CheckPermission("menu:read"), menuHandler.ListMenus)
		menus.GET("/tree", rbac.CheckPermission("menu:read:tree"), menuHandler.GetMenuTree)
		menus.GET("/:id", rbac.CheckPermission("menu:read:detail"), menuHandler.GetMenu)
		menus.PUT("/:id", menuHandler.UpdateMenu)
		menus.DELETE("/:id", menuHandler.DeleteMenu)
		menus.POST("/:id/resources", menuHandler.BindResource)
		menus.DELETE("/:id/resources/:resourceId", menuHandler.UnbindResource)
		menus.GET("/:id/resources", menuHandler.GetMenuResources)
	}
}
