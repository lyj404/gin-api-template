package route

import (
	"github.com/gin-gonic/gin"
	"github.com/lyj404/gin-api-template/api/handler"
	"github.com/lyj404/gin-api-template/api/middleware"
)

func NewUserManagementRouter(h *handler.UserManagementHandler, rbac *middleware.RBACMiddleware, group *gin.RouterGroup) {
	group.GET("/users", rbac.CheckPermission("user:read"), h.ListUsers)
	group.GET("/users/:id", rbac.CheckPermission("user:read:detail"), h.GetUser)
	group.POST("/users", h.CreateUser)
	group.PUT("/users/:id", h.UpdateUser)
	group.DELETE("/users/:id", h.DeleteUser)
}
