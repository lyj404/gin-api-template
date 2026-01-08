package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lyj404/gin-api-template/domain/result"
	"github.com/lyj404/gin-api-template/domain/services"
)

type RBACMiddleware struct {
	permissionService services.PermissionService
}

func NewRBACMiddleware(permissionService services.PermissionService) *RBACMiddleware {
	return &RBACMiddleware{
		permissionService: permissionService,
	}
}

// CheckPermission 检查资源权限中间件
func (m *RBACMiddleware) CheckPermission(resource string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			result.ErrorResponse(c, http.StatusUnauthorized, "未授权")
			c.Abort()
			return
		}

		hasPermission, err := m.permissionService.CheckPermission(userID.(uint), resource, c.Request.Method)
		if err != nil {
			result.ErrorResponse(c, http.StatusInternalServerError, "权限检查失败")
			c.Abort()
			return
		}

		if !hasPermission {
			result.ErrorResponse(c, http.StatusForbidden, "权限不足")
			c.Abort()
			return
		}

		c.Next()
	}
}

// CheckEntityPermission 检查实体权限中间件
func (m *RBACMiddleware) CheckEntityPermission(entityType string, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			result.ErrorResponse(c, http.StatusUnauthorized, "未授权")
			c.Abort()
			return
		}

		entityID := c.Param("id")
		if entityID == "" {
			result.ErrorResponse(c, http.StatusBadRequest, "缺少实体ID")
			c.Abort()
			return
		}

		hasPermission, err := m.permissionService.CheckEntityPermission(userID.(uint), entityType, 0, action)
		if err != nil {
			result.ErrorResponse(c, http.StatusInternalServerError, "权限检查失败")
			c.Abort()
			return
		}

		if !hasPermission {
			result.ErrorResponse(c, http.StatusForbidden, "权限不足")
			c.Abort()
			return
		}

		c.Next()
	}
}
