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
	"github.com/lyj404/gin-api-template/pkg/pagination"
)

// RoleHandler 角色处理器，处理角色相关的HTTP请求
type RoleHandler struct {
	roleService services.RoleService
}

// NewRoleHandler 创建角色处理器实例
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

// ListRoles 获取角色列表（分页）
// @Summary 获取角色列表
// @Description 获取角色列表（支持分页、搜索、排序）
// @Tags 角色
// @Produce json
// @Param page query int false "页码，默认1"
// @Param page_size query int false "每页数量，默认10，最大100"
// @Param keyword query string false "搜索关键词（搜索角色名称）"
// @Param order_by query string false "排序字段"
// @Param sort query string false "排序方式：asc/desc"
// @Success 200 {object} result.ResponseResult[dto.PaginationResponse] "获取成功"
// @Router /roles [get]
func (h *RoleHandler) ListRoles(c *gin.Context) {
	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		result.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	req.SetDefaults()

	var roles []entity.Role
	builder := pagination.NewPaginationBuilder(global.G_DB).
		Model(&entity.Role{}).
		SetPage(req.Page).
		SetPageSize(req.PageSize).
		OrderBy(req.OrderBy + " " + req.Sort)

	// 如果有关键词搜索，添加搜索条件
	if req.Keyword != "" {
		builder = builder.Where("name LIKE ?", "%"+req.Keyword+"%")
	}

	paginationResult, err := builder.Build(&roles)
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

	result.SuccessResponse(c, "获取角色列表成功", dto.NewPaginationResponse(
		paginationResult.Page,
		paginationResult.PageSize,
		paginationResult.Total,
		responses,
	))
}
