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

	operatorID := c.GetUint64("user_id")
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
		G_MODEL:     global.G_MODEL{ID: uint64(id)},
		Name:        request.Name,
		Description: request.Description,
	}

	operatorID := c.GetUint64("user_id")
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

	operatorID := c.GetUint64("user_id")
	if err := h.roleService.DeleteRole(uint64(id), operatorID); err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.SimpleSuccessResponse(c, "角色删除成功")
}

// GetRole 获取角色详情
// @Summary 获取角色详情
// @Description 根据ID获取角色详情（包含绑定的资源列表和菜单列表）
// @Tags 角色
// @Produce json
// @Param id path int true "角色ID"
// @Success 200 {object} result.ResponseResult[dto.RoleDetailResponse] "获取成功"
// @Failure 404 {object} result.ResponseResult[string] "角色不存在"
// @Router /roles/:id [get]
func (h *RoleHandler) GetRole(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	role, err := h.roleService.GetRoleByID(uint64(id))
	if err != nil {
		result.ErrorResponse(c, http.StatusNotFound, "角色不存在")
		return
	}

	resources, _ := h.roleService.GetRoleResources(uint64(id))
	resourceResponses := make([]dto.RoleResourceResponse, len(resources))
	for i, rr := range resources {
		resourceResponses[i] = dto.RoleResourceResponse{
			ID:         rr.ID,
			RoleID:     rr.RoleID,
			ResourceID: rr.ResourceID,
			IsRead:     rr.IsRead,
			IsWrite:    rr.IsWrite,
			Resource: &dto.ResourceBriefResponse{
				ID:          rr.Resource.ID,
				Name:        rr.Resource.Name,
				Type:        rr.Resource.Type,
				Pattern:     rr.Resource.Pattern,
				Method:      rr.Resource.Method,
				Entity:      rr.Resource.Entity,
				Action:      rr.Resource.Action,
				Description: rr.Resource.Description,
			},
		}
	}

	roleMenus, _ := h.roleService.GetRoleMenus(uint64(id))
	menuResponses := make([]dto.RoleMenuResponse, len(roleMenus))
	for i, rm := range roleMenus {
		menuResponses[i] = dto.RoleMenuResponse{
			ID:     rm.ID,
			RoleID: rm.RoleID,
			MenuID: rm.MenuID,
			Menu: &dto.MenuBriefResponse{
				ID:   rm.Menu.ID,
				Name: rm.Menu.Name,
				Path: rm.Menu.Path,
				Icon: rm.Menu.Icon,
			},
		}
	}

	response := dto.RoleDetailResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		IsSystem:    role.IsSystem,
		Resources:   resourceResponses,
		Menus:       menuResponses,
	}

	result.SuccessResponse(c, "获取角色成功", &response)
}

// BindResource 角色绑定资源
// @Summary 角色绑定资源
// @Description 为角色绑定一个资源
// @Tags 角色
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Param request body dto.BindRoleResourceRequest true "绑定资源请求"
// @Success 200 {object} result.ResponseResult[string] "绑定成功"
// @Failure 400 {object} result.ResponseResult[string] "请求参数错误"
// @Failure 500 {object} result.ResponseResult[string] "服务器内部错误"
// @Router /roles/:id/resources [post]
func (h *RoleHandler) BindResource(c *gin.Context) {
	roleID, _ := strconv.Atoi(c.Param("id"))

	var req dto.BindRoleResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	operatorID := c.GetUint64("user_id")
	if err := h.roleService.BindResource(uint64(roleID), req.ResourceID, req.IsWrite, operatorID); err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.SimpleSuccessResponse(c, "资源绑定成功")
}

// UnbindResource 角色解绑资源
// @Summary 角色解绑资源
// @Description 移除角色绑定的资源
// @Tags 角色
// @Produce json
// @Param id path int true "角色ID"
// @Param resourceId path int true "资源ID"
// @Success 200 {object} result.ResponseResult[string] "解绑成功"
// @Failure 500 {object} result.ResponseResult[string] "服务器内部错误"
// @Router /roles/:id/resources/:resourceId [delete]
func (h *RoleHandler) UnbindResource(c *gin.Context) {
	roleID, _ := strconv.Atoi(c.Param("id"))
	resourceID, _ := strconv.Atoi(c.Param("resourceId"))

	operatorID := c.GetUint64("user_id")
	if err := h.roleService.UnbindResource(uint64(roleID), uint64(resourceID), operatorID); err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.SimpleSuccessResponse(c, "资源解绑成功")
}

// GetRoleResources 获取角色绑定的资源列表
// @Summary 获取角色绑定的资源列表
// @Description 获取指定角色绑定的所有资源
// @Tags 角色
// @Produce json
// @Param id path int true "角色ID"
// @Success 200 {object} result.ResponseResult[[]dto.RoleResourceResponse] "获取成功"
// @Failure 500 {object} result.ResponseResult[string] "服务器内部错误"
// @Router /roles/:id/resources [get]
func (h *RoleHandler) GetRoleResources(c *gin.Context) {
	roleID, _ := strconv.Atoi(c.Param("id"))

	resources, err := h.roleService.GetRoleResources(uint64(roleID))
	if err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	responses := make([]dto.RoleResourceResponse, len(resources))
	for i, rr := range resources {
		responses[i] = dto.RoleResourceResponse{
			ID:         rr.ID,
			RoleID:     rr.RoleID,
			ResourceID: rr.ResourceID,
			IsRead:     rr.IsRead,
			IsWrite:    rr.IsWrite,
			Resource: &dto.ResourceBriefResponse{
				ID:          rr.Resource.ID,
				Name:        rr.Resource.Name,
				Type:        rr.Resource.Type,
				Pattern:     rr.Resource.Pattern,
				Method:      rr.Resource.Method,
				Entity:      rr.Resource.Entity,
				Action:      rr.Resource.Action,
				Description: rr.Resource.Description,
			},
		}
	}

	result.SuccessResponse(c, "获取角色资源成功", &responses)
}

// BindMenu 角色绑定菜单
// @Summary 角色绑定菜单
// @Description 为角色绑定一个菜单
// @Tags 角色
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Param request body dto.BindRoleMenuRequest true "绑定菜单请求"
// @Success 200 {object} result.ResponseResult[string] "绑定成功"
// @Failure 400 {object} result.ResponseResult[string] "请求参数错误"
// @Failure 500 {object} result.ResponseResult[string] "服务器内部错误"
// @Router /roles/:id/menus [post]
func (h *RoleHandler) BindMenu(c *gin.Context) {
	roleID, _ := strconv.Atoi(c.Param("id"))

	var req dto.BindRoleMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	operatorID := c.GetUint64("user_id")
	if err := h.roleService.BindMenu(uint64(roleID), req.MenuID, operatorID); err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.SimpleSuccessResponse(c, "菜单绑定成功")
}

// UnbindMenu 角色解绑菜单
// @Summary 角色解绑菜单
// @Description 移除角色绑定的菜单
// @Tags 角色
// @Produce json
// @Param id path int true "角色ID"
// @Param menuId path int true "菜单ID"
// @Success 200 {object} result.ResponseResult[string] "解绑成功"
// @Failure 500 {object} result.ResponseResult[string] "服务器内部错误"
// @Router /roles/:id/menus/:menuId [delete]
func (h *RoleHandler) UnbindMenu(c *gin.Context) {
	roleID, _ := strconv.Atoi(c.Param("id"))
	menuID, _ := strconv.Atoi(c.Param("menuId"))

	operatorID := c.GetUint64("user_id")
	if err := h.roleService.UnbindMenu(uint64(roleID), uint64(menuID), operatorID); err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.SimpleSuccessResponse(c, "菜单解绑成功")
}

// GetRoleMenus 获取角色绑定的菜单列表
// @Summary 获取角色绑定的菜单列表
// @Description 获取指定角色绑定的所有菜单
// @Tags 角色
// @Produce json
// @Param id path int true "角色ID"
// @Success 200 {object} result.ResponseResult[[]dto.RoleMenuResponse] "获取成功"
// @Failure 500 {object} result.ResponseResult[string] "服务器内部错误"
// @Router /roles/:id/menus [get]
func (h *RoleHandler) GetRoleMenus(c *gin.Context) {
	roleID, _ := strconv.Atoi(c.Param("id"))

	roleMenus, err := h.roleService.GetRoleMenus(uint64(roleID))
	if err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	responses := make([]dto.RoleMenuResponse, len(roleMenus))
	for i, rm := range roleMenus {
		responses[i] = dto.RoleMenuResponse{
			ID:     rm.ID,
			RoleID: rm.RoleID,
			MenuID: rm.MenuID,
			Menu: &dto.MenuBriefResponse{
				ID:   rm.Menu.ID,
				Name: rm.Menu.Name,
				Path: rm.Menu.Path,
				Icon: rm.Menu.Icon,
			},
		}
	}

	result.SuccessResponse(c, "获取角色菜单成功", &responses)
}

// ListRoles 获取角色列表（分页）
// @Summary 获取角色列表
// @Description 获取角色列表（支持分页、搜索、排序）。系统管理员可查看所有角色，非系统角色用户只能查看非系统角色。
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

	orderBy := req.OrderBy
	if orderBy == "" {
		orderBy = "id"
	}
	orderBy += " " + req.Sort

	operatorID := c.GetUint64("user_id")

	var roleIDs []uint64
	global.G_DB.Model(&entity.UserRole{}).Select("role_id").Where("user_id = ?", operatorID).Find(&roleIDs)
	hasSystemRole := false
	if len(roleIDs) > 0 {
		var cnt int64
		global.G_DB.Model(&entity.Role{}).Where("id IN ? AND is_system = ?", roleIDs, true).Count(&cnt)
		hasSystemRole = cnt > 0
	}

	var roles []entity.Role
	builder := pagination.NewPaginationBuilder(global.G_DB).
		Model(&entity.Role{}).
		SetPage(req.Page).
		SetPageSize(req.PageSize).
		OrderBy(orderBy)

	if req.Keyword != "" {
		builder = builder.Where("name LIKE ?", "%"+req.Keyword+"%")
	}

	// 非系统角色只能看到非系统角色
	if !hasSystemRole {
		builder = builder.Where("is_system = ?", false)
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
