package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lyj404/gin-api-template/domain/dto"
	"github.com/lyj404/gin-api-template/domain/entity"
	"github.com/lyj404/gin-api-template/domain/result"
	"github.com/lyj404/gin-api-template/domain/services"
	"github.com/lyj404/gin-api-template/global"
)

type RoleHandler struct {
	roleService services.RoleService
}

func NewRoleHandler(roleService services.RoleService) *RoleHandler {
	return &RoleHandler{
		roleService: roleService,
	}
}

// CreateRole 创建角色
// @Summary 创建角色
// @Description 创建新角色
// @Tags 角色
// @Accept json
// @Produce json
// @Param request body dto.CreateRoleRequest true "角色信息"
// @Success 200 {object} result.ResponseResult[dto.RoleResponse] "创建成功"
// @Failure 400 {object} result.ResponseResult[string] "请求参数错误"
// @Failure 500 {object} result.ResponseResult[string] "服务器内部错误"
// @Router /roles [post]
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var request dto.CreateRoleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		result.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	role := &entity.Role{
		Name:        request.Name,
		Description: request.Description,
		IsSystem:    false,
	}

	operatorID := c.GetUint("user_id")
	if err := h.roleService.CreateRole(role, operatorID); err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := dto.RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		IsSystem:    role.IsSystem,
	}

	result.SuccessResponse(c, "角色创建成功", &response)
}

// UpdateRole 更新角色
// @Summary 更新角色
// @Description 更新角色信息
// @Tags 角色
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Param request body dto.UpdateRoleRequest true "角色信息"
// @Success 200 {object} result.ResponseResult[dto.RoleResponse] "更新成功"
// @Failure 400 {object} result.ResponseResult[string] "请求参数错误"
// @Failure 404 {object} result.ResponseResult[string] "角色不存在"
// @Failure 500 {object} result.ResponseResult[string] "服务器内部错误"
// @Router /roles/:id [put]
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var request dto.UpdateRoleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		result.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	role := &entity.Role{
		G_MODEL:     global.G_MODEL{ID: uint(id)},
		Name:        request.Name,
		Description: request.Description,
	}

	operatorID := c.GetUint("user_id")
	if err := h.roleService.UpdateRole(role, operatorID); err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := dto.RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		IsSystem:    role.IsSystem,
	}

	result.SuccessResponse(c, "角色更新成功", &response)
}

// DeleteRole 删除角色
// @Summary 删除角色
// @Description 删除角色
// @Tags 角色
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Success 200 {object} result.ResponseResult[string] "删除成功"
// @Failure 400 {object} result.ResponseResult[string] "请求参数错误"
// @Failure 500 {object} result.ResponseResult[string] "服务器内部错误"
// @Router /roles/:id [delete]
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	operatorID := c.GetUint("user_id")
	if err := h.roleService.DeleteRole(uint(id), operatorID); err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.SimpleSuccessResponse(c, "角色删除成功")
}

// GetRole 获取角色详情
// @Summary 获取角色详情
// @Description 根据ID获取角色详情
// @Tags 角色
// @Produce json
// @Param id path int true "角色ID"
// @Success 200 {object} result.ResponseResult[dto.RoleResponse] "获取成功"
// @Failure 404 {object} result.ResponseResult[string] "角色不存在"
// @Router /roles/:id [get]
func (h *RoleHandler) GetRole(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	role, err := h.roleService.GetRoleByID(uint(id))
	if err != nil {
		result.ErrorResponse(c, http.StatusNotFound, "角色不存在")
		return
	}

	response := dto.RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		IsSystem:    role.IsSystem,
	}

	result.SuccessResponse(c, "获取角色成功", &response)
}

// ListRoles 获取角色列表
// @Summary 获取角色列表
// @Description 获取所有角色
// @Tags 角色
// @Produce json
// @Success 200 {object} result.ResponseResult[[]dto.RoleResponse] "获取成功"
// @Router /roles [get]
func (h *RoleHandler) ListRoles(c *gin.Context) {
	roles, err := h.roleService.GetAllRoles()
	if err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	responses := make([]dto.RoleResponse, len(roles))
	for i, role := range roles {
		responses[i] = dto.RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			IsSystem:    role.IsSystem,
		}
	}

	result.SuccessResponse(c, "获取角色列表成功", &responses)
}
