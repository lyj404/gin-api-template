package route

import (
	"github.com/gin-gonic/gin"
	"github.com/lyj404/gin-api-template/api/handler"
)

// NewMenuRouter 注册菜单路由
func NewMenuRouter(menuHandler *handler.MenuHandler, group *gin.RouterGroup) {
	menus := group.Group("/menus")
	{
		menus.POST("", menuHandler.CreateMenu)
		menus.GET("", menuHandler.ListMenus)
		menus.GET("/tree", menuHandler.GetMenuTree)
		menus.GET("/:id", menuHandler.GetMenu)
		menus.PUT("/:id", menuHandler.UpdateMenu)
		menus.DELETE("/:id", menuHandler.DeleteMenu)
	}
}