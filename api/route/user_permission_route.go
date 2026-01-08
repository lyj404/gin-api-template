package route

import (
	"github.com/gin-gonic/gin"
	"github.com/lyj404/gin-api-template/api/handler"
)

func NewUserPermissionRouter(userPermHdlr *handler.UserPermissionHandler, group *gin.RouterGroup) {
	group.GET("/user/permissions", userPermHdlr.GetUserPermissions)
}
