package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lyj404/gin-api-template/domain/result"
	"github.com/lyj404/gin-api-template/domain/services"
)

// RBACMiddleware 基于角色的访问控制中间件
type RBACMiddleware struct {
	permissionService services.PermissionService
}

// NewRBACMiddleware 创建RBAC中间件实例
func NewRBACMiddleware(permissionService services.PermissionService) *RBACMiddleware {
	return &RBACMiddleware{
		permissionService: permissionService,
	}
}

// CheckPermission 检查资源权限中间件
// 该中间件验证用户是否有权限访问指定的资源路径
func (m *RBACMiddleware) CheckPermission(resource string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			result.ErrorResponse(c, http.StatusUnauthorized, "未授权")
			c.Abort()
			return
		}

		hasPermission, err := m.permissionService.CheckPermission(userID.(uint64), resource, c.Request.Method)
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
// 该中间件验证用户是否有权限对指定实体执行特定操作
func (m *RBACMiddleware) CheckEntityPermission(entityType string, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			result.ErrorResponse(c, http.StatusUnauthorized, "未授权")
			c.Abort()
			return
		}

		// 从URL参数中解析实体ID
		entityIDStr := c.Param("id")
		if entityIDStr == "" {
			result.ErrorResponse(c, http.StatusBadRequest, "缺少实体ID")
			c.Abort()
			return
		}

		entityID, err := strconv.ParseUint(entityIDStr, 10, 64)
		if err != nil {
			result.ErrorResponse(c, http.StatusBadRequest, "无效的实体ID")
			c.Abort()
			return
		}

		hasPermission, err := m.permissionService.CheckEntityPermission(userID.(uint64), entityType, uint64(entityID), action)
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
