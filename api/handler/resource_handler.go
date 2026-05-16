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

type ResourceHandler struct {
	resourceService services.ResourceService
}

func NewResourceHandler(resourceService services.ResourceService) *ResourceHandler {
	return &ResourceHandler{
		resourceService: resourceService,
	}
}

func (h *ResourceHandler) CreateResource(c *gin.Context) {
	var request dto.CreateResourceRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		result.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	resource := &entity.Resource{
		Name:        request.Name,
		Type:        request.Type,
		Pattern:     request.Pattern,
		Method:      request.Method,
		Entity:      request.Entity,
		Action:      request.Action,
		Description: request.Description,
	}

	operatorID := c.GetUint64("user_id")
	if err := h.resourceService.CreateResource(resource, operatorID); err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := dto.ResourceListResponse{
		ID:          resource.ID,
		Name:        resource.Name,
		Type:        resource.Type,
		Pattern:     resource.Pattern,
		Method:      resource.Method,
		Entity:      resource.Entity,
		Action:      resource.Action,
		Description: resource.Description,
	}

	result.SuccessResponse(c, "资源创建成功", &response)
}

func (h *ResourceHandler) UpdateResource(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var request dto.UpdateResourceRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		result.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	resource := &entity.Resource{
		G_MODEL:     global.G_MODEL{ID: id},
		Name:        request.Name,
		Type:        request.Type,
		Pattern:     request.Pattern,
		Method:      request.Method,
		Entity:      request.Entity,
		Action:      request.Action,
		Description: request.Description,
	}

	operatorID := c.GetUint64("user_id")
	if err := h.resourceService.UpdateResource(resource, operatorID); err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := dto.ResourceListResponse{
		ID:          resource.ID,
		Name:        resource.Name,
		Type:        resource.Type,
		Pattern:     resource.Pattern,
		Method:      resource.Method,
		Entity:      resource.Entity,
		Action:      resource.Action,
		Description: resource.Description,
	}

	result.SuccessResponse(c, "资源更新成功", &response)
}

func (h *ResourceHandler) DeleteResource(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	operatorID := c.GetUint64("user_id")
	if err := h.resourceService.DeleteResource(id, operatorID); err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.SimpleSuccessResponse(c, "资源删除成功")
}

func (h *ResourceHandler) GetResource(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	resource, err := h.resourceService.GetResourceByID(id)
	if err != nil {
		result.ErrorResponse(c, http.StatusNotFound, "资源不存在")
		return
	}

	response := dto.ResourceListResponse{
		ID:          resource.ID,
		Name:        resource.Name,
		Type:        resource.Type,
		Pattern:     resource.Pattern,
		Method:      resource.Method,
		Entity:      resource.Entity,
		Action:      resource.Action,
		Description: resource.Description,
	}

	result.SuccessResponse(c, "获取资源成功", &response)
}

func (h *ResourceHandler) ListResources(c *gin.Context) {
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

	var resources []entity.Resource
	builder := pagination.NewPaginationBuilder(global.G_DB).
		Model(&entity.Resource{}).
		SetPage(req.Page).
		SetPageSize(req.PageSize).
		OrderBy(orderBy)

	if req.Keyword != "" {
		builder = builder.Where("name LIKE ?", "%"+req.Keyword+"%")
	}

	paginationResult, err := builder.Build(&resources)
	if err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	responses := make([]dto.ResourceListResponse, len(resources))
	for i, r := range resources {
		responses[i] = dto.ResourceListResponse{
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

	result.SuccessResponse(c, "获取资源列表成功", dto.NewPaginationResponse(
		paginationResult.Page,
		paginationResult.PageSize,
		paginationResult.Total,
		responses,
	))
}
