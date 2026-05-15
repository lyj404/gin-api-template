package route

import (
	"github.com/gin-gonic/gin"
	"github.com/lyj404/gin-api-template/api/handler"
)

func NewMenuRouter(menuHandler *handler.MenuHandler, group *gin.RouterGroup) {
	menus := group.Group("/menus")
	{
		menus.POST("", menuHandler.CreateMenu)
		menus.GET("", menuHandler.ListMenus)
		menus.GET("/tree", menuHandler.GetMenuTree)
		menus.GET("/:id", menuHandler.GetMenu)
		menus.PUT("/:id", menuHandler.UpdateMenu)
		menus.DELETE("/:id", menuHandler.DeleteMenu)
		menus.POST("/:id/resources", menuHandler.BindResource)
		menus.DELETE("/:id/resources/:resourceId", menuHandler.UnbindResource)
		menus.GET("/:id/resources", menuHandler.GetMenuResources)
	}
}
