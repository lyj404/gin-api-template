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

type MenuHandler struct {
	menuService services.MenuService
}

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
		Name:      request.Name,
		ParentID:  request.ParentID,
		Path:      request.Path,
		Icon:      request.Icon,
		OrderNum:  request.OrderNum,
		IsVisible: request.IsVisible,
		Status:    "enabled",
	}

	operatorID := c.GetUint64("user_id")
	if err := h.menuService.CreateMenu(menu, operatorID); err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := dto.MenuResponse{
		ID:        menu.ID,
		Name:      menu.Name,
		ParentID:  menu.ParentID,
		Path:      menu.Path,
		Icon:      menu.Icon,
		OrderNum:  menu.OrderNum,
		IsVisible: menu.IsVisible,
		Status:    menu.Status,
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
		G_MODEL:   global.G_MODEL{ID: uint64(id)},
		Name:      request.Name,
		ParentID:  request.ParentID,
		Path:      request.Path,
		Icon:      request.Icon,
		OrderNum:  request.OrderNum,
		Status:    request.Status,
	}

	if request.IsVisible != nil {
		menu.IsVisible = *request.IsVisible
	}

	operatorID := c.GetUint64("user_id")
	if err := h.menuService.UpdateMenu(menu, operatorID); err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := dto.MenuResponse{
		ID:        menu.ID,
		Name:      menu.Name,
		ParentID:  menu.ParentID,
		Path:      menu.Path,
		Icon:      menu.Icon,
		OrderNum:  menu.OrderNum,
		IsVisible: menu.IsVisible,
		Status:    menu.Status,
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

	operatorID := c.GetUint64("user_id")
	if err := h.menuService.DeleteMenu(uint64(id), operatorID); err != nil {
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

	menu, err := h.menuService.GetMenuByID(uint64(id))
	if err != nil {
		result.ErrorResponse(c, http.StatusNotFound, "菜单不存在")
		return
	}

	response := dto.MenuResponse{
		ID:        menu.ID,
		Name:      menu.Name,
		ParentID:  menu.ParentID,
		Path:      menu.Path,
		Icon:      menu.Icon,
		OrderNum:  menu.OrderNum,
		IsVisible: menu.IsVisible,
		Status:    menu.Status,
		Resources: toResourceBriefResponses(menu.Resources),
	}

	result.SuccessResponse(c, "获取菜单成功", &response)
}

// ListMenus 获取菜单列表（分页）
// @Summary 获取菜单列表
// @Description 获取菜单列表（支持分页、搜索、排序）
// @Tags 菜单
// @Produce json
// @Param page query int false "页码，默认1"
// @Param page_size query int false "每页数量，默认10，最大100"
// @Param keyword query string false "搜索关键词（搜索菜单名称）"
// @Param order_by query string false "排序字段"
// @Param sort query string false "排序方式：asc/desc"
// @Success 200 {object} result.ResponseResult[dto.PaginationResponse] "获取成功"
// @Router /menus [get]
func (h *MenuHandler) ListMenus(c *gin.Context) {
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

	var menus []entity.Menu
	builder := pagination.NewPaginationBuilder(global.G_DB).
		Model(&entity.Menu{}).
		SetPage(req.Page).
		SetPageSize(req.PageSize).
		OrderBy(orderBy)

	if req.Keyword != "" {
		builder = builder.Where("name LIKE ?", "%"+req.Keyword+"%")
	}

	paginationResult, err := builder.Build(&menus)
	if err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	responses := make([]dto.MenuListResponse, len(menus))
	for i, menu := range menus {
		responses[i] = dto.MenuListResponse{
			ID:        menu.ID,
			Name:      menu.Name,
			ParentID:  menu.ParentID,
			Path:      menu.Path,
			Icon:      menu.Icon,
			OrderNum:  menu.OrderNum,
			IsVisible: menu.IsVisible,
			Status:    menu.Status,
		}
	}

	result.SuccessResponse(c, "获取菜单列表成功", dto.NewPaginationResponse(
		paginationResult.Page,
		paginationResult.PageSize,
		paginationResult.Total,
		responses,
	))
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

// BindResource 菜单绑定资源
// @Summary 菜单绑定资源
// @Description 为菜单绑定一个资源
// @Tags 菜单
// @Accept json
// @Produce json
// @Param id path int true "菜单ID"
// @Param request body dto.BindMenuResourceRequest true "绑定资源请求"
// @Success 200 {object} result.ResponseResult[string] "绑定成功"
// @Failure 400 {object} result.ResponseResult[string] "请求参数错误"
// @Failure 500 {object} result.ResponseResult[string] "服务器内部错误"
// @Router /menus/:id/resources [post]
func (h *MenuHandler) BindResource(c *gin.Context) {
	menuID, _ := strconv.Atoi(c.Param("id"))

	var req dto.BindMenuResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	operatorID := c.GetUint64("user_id")
	if err := h.menuService.BindResource(uint64(menuID), req.ResourceID, operatorID); err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.SimpleSuccessResponse(c, "资源绑定成功")
}

// UnbindResource 菜单解绑资源
// @Summary 菜单解绑资源
// @Description 移除菜单绑定的资源
// @Tags 菜单
// @Produce json
// @Param id path int true "菜单ID"
// @Param resourceId path int true "资源ID"
// @Success 200 {object} result.ResponseResult[string] "解绑成功"
// @Failure 500 {object} result.ResponseResult[string] "服务器内部错误"
// @Router /menus/:id/resources/:resourceId [delete]
func (h *MenuHandler) UnbindResource(c *gin.Context) {
	menuID, _ := strconv.Atoi(c.Param("id"))
	resourceID, _ := strconv.Atoi(c.Param("resourceId"))

	operatorID := c.GetUint64("user_id")
	if err := h.menuService.UnbindResource(uint64(menuID), uint64(resourceID), operatorID); err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.SimpleSuccessResponse(c, "资源解绑成功")
}

// GetMenuResources 获取菜单绑定的资源列表
// @Summary 获取菜单绑定的资源列表
// @Description 获取指定菜单绑定的所有资源
// @Tags 菜单
// @Produce json
// @Param id path int true "菜单ID"
// @Success 200 {object} result.ResponseResult[[]dto.MenuResourceResponse] "获取成功"
// @Failure 500 {object} result.ResponseResult[string] "服务器内部错误"
// @Router /menus/:id/resources [get]
func (h *MenuHandler) GetMenuResources(c *gin.Context) {
	menuID, _ := strconv.Atoi(c.Param("id"))

	resources, err := h.menuService.GetMenuResources(uint64(menuID))
	if err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	responses := make([]dto.MenuResourceResponse, len(resources))
	for i, mr := range resources {
		responses[i] = dto.MenuResourceResponse{
			ID:         mr.ID,
			MenuID:     mr.MenuID,
			ResourceID: mr.ResourceID,
			Resource: &dto.ResourceBriefResponse{
				ID:          mr.Resource.ID,
				Name:        mr.Resource.Name,
				Type:        mr.Resource.Type,
				Pattern:     mr.Resource.Pattern,
				Method:      mr.Resource.Method,
				Entity:      mr.Resource.Entity,
				Action:      mr.Resource.Action,
				Description: mr.Resource.Description,
			},
		}
	}

	result.SuccessResponse(c, "获取菜单资源成功", &responses)
}

func convertMenuTreeToResponse(menus []entity.Menu) []dto.MenuResponse {
	responses := make([]dto.MenuResponse, 0, len(menus))
	for _, menu := range menus {
		response := dto.MenuResponse{
			ID:        menu.ID,
			Name:      menu.Name,
			ParentID:  menu.ParentID,
			Path:      menu.Path,
			Icon:      menu.Icon,
			OrderNum:  menu.OrderNum,
			IsVisible: menu.IsVisible,
			Status:    menu.Status,
			Resources: toResourceBriefResponses(menu.Resources),
		}
		if len(menu.Children) > 0 {
			response.Children = convertMenuTreeToResponse(menu.Children)
		}
		responses = append(responses, response)
	}
	return responses
}

func toResourceBriefResponses(resources []entity.Resource) []dto.ResourceBriefResponse {
	if len(resources) == 0 {
		return nil
	}
	resp := make([]dto.ResourceBriefResponse, len(resources))
	for i, r := range resources {
		resp[i] = dto.ResourceBriefResponse{
			ID:          r.ID,
			Name:        r.Name,
			Type:        r.Type,
			Pattern:     r.Pattern,
			Method:      r.Method,
			Entity:      r.Entity,
			Action:      r.Action,
			Description: r.Description,
		}
	}
	return resp
}
