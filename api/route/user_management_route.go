package route

import (
	"github.com/gin-gonic/gin"
	"github.com/lyj404/gin-api-template/api/handler"
)

func NewUserManagementRouter(h *handler.UserManagementHandler, group *gin.RouterGroup) {
	group.GET("/users", h.ListUsers)
	group.GET("/users/:id", h.GetUser)
	group.POST("/users", h.CreateUser)
	group.PUT("/users/:id", h.UpdateUser)
	group.DELETE("/users/:id", h.DeleteUser)
}
