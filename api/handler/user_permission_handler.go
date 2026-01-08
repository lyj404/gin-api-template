package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lyj404/gin-api-template/domain/result"
	domainservices "github.com/lyj404/gin-api-template/domain/services"
)

type UserPermissionHandler struct {
	permissionService domainservices.PermissionService
}

func NewUserPermissionHandler(permissionService domainservices.PermissionService) *UserPermissionHandler {
	return &UserPermissionHandler{
		permissionService: permissionService,
	}
}

// GetUserPermissions 获取当前用户权限
// @Summary 获取当前用户权限
// @Description 获取当前用户的权限列表和组织范围
// @Tags 用户
// @Produce json
// @Success 200 {object} result.ResponseResult[UserPermissionResponse] "获取成功"
// @Failure 401 {object} result.ResponseResult[string] "未授权"
// @Failure 500 {object} result.ResponseResult[string] "服务器内部错误"
// @Router /user/permissions [get]
func (h *UserPermissionHandler) GetUserPermissions(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		result.ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}

	permissions, err := h.permissionService.GetUserPermissions(userID.(uint))
	if err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	orgScopes, err := h.permissionService.GetUserOrgScope(userID.(uint))
	if err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := UserPermissionResponse{
		Permissions: permissions,
		OrgScopes:   orgScopes,
	}

	result.SuccessResponse(c, "获取用户权限成功", &response)
}

type UserPermissionResponse struct {
	Permissions []domainservices.PermissionInfo `json:"permissions"`
	OrgScopes   []domainservices.OrgScopeInfo   `json:"org_scopes"`
}
