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

// MenuHandler 菜单处理器，处理菜单相关的HTTP请求
type MenuHandler struct {
	menuService services.MenuService
}

// NewMenuHandler 创建菜单处理器实例
func NewMenuHandler(menuService services.MenuService) *MenuHandler {
	return &MenuHandler{
		menuService: menuService,
	}
}

// CreateMenu 创建菜单
// @Summary 创建菜单
// @Description 创建新菜单
// @Tags 菜单
// @Accept json
// @Produce json
// @Param request body dto.CreateMenuRequest true "菜单信息"
// @Success 200 {object} result.ResponseResult[dto.MenuResponse] "创建成功"
// @Failure 400 {object} result.ResponseResult[string] "请求参数错误"
// @Failure 500 {object} result.ResponseResult[string] "服务器内部错误"
// @Router /menus [post]
func (h *MenuHandler) CreateMenu(c *gin.Context) {
	var request dto.CreateMenuRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		result.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	menu := &entity.Menu{
		Name:       request.Name,
		ParentID:   request.ParentID,
		Path:       request.Path,
		Component:  request.Component,
		Icon:       request.Icon,
		OrderNum:   request.OrderNum,
		ResourceID: request.ResourceID,
		IsVisible:  request.IsVisible,
		Status:     "enabled",
	}

	operatorID := c.GetUint("user_id")
	if err := h.menuService.CreateMenu(menu, operatorID); err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := dto.MenuResponse{
		ID:         menu.ID,
		Name:       menu.Name,
		ParentID:   menu.ParentID,
		Path:       menu.Path,
		Component:  menu.Component,
		Icon:       menu.Icon,
		OrderNum:   menu.OrderNum,
		ResourceID: menu.ResourceID,
		IsVisible:  menu.IsVisible,
		Status:     menu.Status,
	}

	result.SuccessResponse(c, "菜单创建成功", &response)
}

// UpdateMenu 更新菜单
// @Summary 更新菜单
// @Description 更新菜单信息
// @Tags 菜单
// @Accept json
// @Produce json
// @Param id path int true "菜单ID"
// @Param request body dto.UpdateMenuRequest true "菜单信息"
// @Success 200 {object} result.ResponseResult[dto.MenuResponse] "更新成功"
// @Failure 400 {object} result.ResponseResult[string] "请求参数错误"
// @Failure 404 {object} result.ResponseResult[string] "菜单不存在"
// @Failure 500 {object} result.ResponseResult[string] "服务器内部错误"
// @Router /menus/:id [put]
func (h *MenuHandler) UpdateMenu(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var request dto.UpdateMenuRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		result.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	menu := &entity.Menu{
		G_MODEL:    global.G_MODEL{ID: uint(id)},
		Name:       request.Name,
		ParentID:   request.ParentID,
		Path:       request.Path,
		Component:  request.Component,
		Icon:       request.Icon,
		OrderNum:   request.OrderNum,
		ResourceID: request.ResourceID,
		Status:     request.Status,
	}

	// 如果 IsVisible 被传入，更新它
	if request.IsVisible != nil {
		menu.IsVisible = *request.IsVisible
	}

	operatorID := c.GetUint("user_id")
	if err := h.menuService.UpdateMenu(menu, operatorID); err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := dto.MenuResponse{
		ID:         menu.ID,
		Name:       menu.Name,
		ParentID:   menu.ParentID,
		Path:       menu.Path,
		Component:  menu.Component,
		Icon:       menu.Icon,
		OrderNum:   menu.OrderNum,
		ResourceID: menu.ResourceID,
		IsVisible:  menu.IsVisible,
		Status:     menu.Status,
	}

	result.SuccessResponse(c, "菜单更新成功", &response)
}

// DeleteMenu 删除菜单
// @Summary 删除菜单
// @Description 删除菜单
// @Tags 菜单
// @Accept json
// @Produce json
// @Param id path int true "菜单ID"
// @Success 200 {object} result.ResponseResult[string] "删除成功"
// @Failure 400 {object} result.ResponseResult[string] "请求参数错误"
// @Failure 500 {object} result.ResponseResult[string] "服务器内部错误"
// @Router /menus/:id [delete]
func (h *MenuHandler) DeleteMenu(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	operatorID := c.GetUint("user_id")
	if err := h.menuService.DeleteMenu(uint(id), operatorID); err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.SimpleSuccessResponse(c, "菜单删除成功")
}

// GetMenu 获取菜单详情
// @Summary 获取菜单详情
// @Description 根据ID获取菜单详情
// @Tags 菜单
// @Produce json
// @Param id path int true "菜单ID"
// @Success 200 {object} result.ResponseResult[dto.MenuResponse] "获取成功"
// @Failure 404 {object} result.ResponseResult[string] "菜单不存在"
// @Router /menus/:id [get]
func (h *MenuHandler) GetMenu(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	menu, err := h.menuService.GetMenuByID(uint(id))
	if err != nil {
		result.ErrorResponse(c, http.StatusNotFound, "菜单不存在")
		return
	}

	response := dto.MenuResponse{
		ID:         menu.ID,
		Name:       menu.Name,
		ParentID:   menu.ParentID,
		Path:       menu.Path,
		Component:  menu.Component,
		Icon:       menu.Icon,
		OrderNum:   menu.OrderNum,
		ResourceID: menu.ResourceID,
		IsVisible:  menu.IsVisible,
		Status:     menu.Status,
	}

	result.SuccessResponse(c, "获取菜单成功", &response)
}

// ListMenus 获取菜单列表
// @Summary 获取菜单列表
// @Description 获取所有菜单（平面结构）
// @Tags 菜单
// @Produce json
// @Success 200 {object} result.ResponseResult[[]dto.MenuListResponse] "获取成功"
// @Router /menus [get]
func (h *MenuHandler) ListMenus(c *gin.Context) {
	menus, err := h.menuService.GetAllMenus()
	if err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	responses := make([]dto.MenuListResponse, len(menus))
	for i, menu := range menus {
		resourceName := ""
		if menu.Resource != nil {
			resourceName = menu.Resource.Name
		}
		responses[i] = dto.MenuListResponse{
			ID:           menu.ID,
			Name:         menu.Name,
			ParentID:     menu.ParentID,
			Path:         menu.Path,
			Component:    menu.Component,
			Icon:         menu.Icon,
			OrderNum:     menu.OrderNum,
			ResourceID:   menu.ResourceID,
			ResourceName: resourceName,
			IsVisible:    menu.IsVisible,
			Status:       menu.Status,
		}
	}

	result.SuccessResponse(c, "获取菜单列表成功", &responses)
}

// GetMenuTree 获取菜单树
// @Summary 获取菜单树
// @Description 获取菜单树形结构
// @Tags 菜单
// @Produce json
// @Success 200 {object} result.ResponseResult[[]dto.MenuResponse] "获取成功"
// @Router /menus/tree [get]
func (h *MenuHandler) GetMenuTree(c *gin.Context) {
	menus, err := h.menuService.GetMenuTree()
	if err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	responses := convertMenuTreeToResponse(menus)
	result.SuccessResponse(c, "获取菜单树成功", &responses)
}

// convertMenuTreeToResponse 将菜单树转换为响应结构
func convertMenuTreeToResponse(menus []entity.Menu) []dto.MenuResponse {
	responses := make([]dto.MenuResponse, 0, len(menus))
	for _, menu := range menus {
		response := dto.MenuResponse{
			ID:         menu.ID,
			Name:       menu.Name,
			ParentID:   menu.ParentID,
			Path:       menu.Path,
			Component:  menu.Component,
			Icon:       menu.Icon,
			OrderNum:   menu.OrderNum,
			ResourceID: menu.ResourceID,
			IsVisible:  menu.IsVisible,
			Status:     menu.Status,
		}
		if len(menu.Children) > 0 {
			response.Children = convertMenuTreeToResponse(menu.Children)
		}
		responses = append(responses, response)
	}
	return responses
}