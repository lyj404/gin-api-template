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

type OrgUnitHandler struct {
	orgService services.OrgUnitService
}

func NewOrgUnitHandler(orgService services.OrgUnitService) *OrgUnitHandler {
	return &OrgUnitHandler{
		orgService: orgService,
	}
}

// CreateOrgUnit 创建组织节点
// @Summary 创建组织节点
// @Description 创建新的组织节点
// @Tags 组织
// @Accept json
// @Produce json
// @Param request body dto.CreateOrgUnitRequest true "组织信息"
// @Success 200 {object} result.ResponseResult[dto.OrgUnitResponse] "创建成功"
// @Failure 400 {object} result.ResponseResult[string] "请求参数错误"
// @Failure 500 {object} result.ResponseResult[string] "服务器内部错误"
// @Router /org-units [post]
func (h *OrgUnitHandler) CreateOrgUnit(c *gin.Context) {
	var request dto.CreateOrgUnitRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		result.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	orgUnit := &entity.OrgUnit{
		Name:     request.Name,
		ParentID: request.ParentID,
	}

	operatorID := c.GetUint("user_id")
	if err := h.orgService.CreateOrgUnit(orgUnit, operatorID); err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := dto.OrgUnitResponse{
		ID:       orgUnit.ID,
		Name:     orgUnit.Name,
		ParentID: orgUnit.ParentID,
		Path:     orgUnit.Path,
		Level:    orgUnit.Level,
	}

	result.SuccessResponse(c, "组织节点创建成功", &response)
}

// UpdateOrgUnit 更新组织节点
// @Summary 更新组织节点
// @Description 更新组织节点信息
// @Tags 组织
// @Accept json
// @Produce json
// @Param id path int true "组织ID"
// @Param request body dto.UpdateOrgUnitRequest true "组织信息"
// @Success 200 {object} result.ResponseResult[dto.OrgUnitResponse] "更新成功"
// @Failure 400 {object} result.ResponseResult[string] "请求参数错误"
// @Failure 404 {object} result.ResponseResult[string] "组织不存在"
// @Failure 500 {object} result.ResponseResult[string] "服务器内部错误"
// @Router /org-units/:id [put]
func (h *OrgUnitHandler) UpdateOrgUnit(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var request dto.UpdateOrgUnitRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		result.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	orgUnit := &entity.OrgUnit{
		G_MODEL:  global.G_MODEL{ID: uint(id)},
		Name:     request.Name,
		ParentID: request.ParentID,
	}

	operatorID := c.GetUint("user_id")
	if err := h.orgService.UpdateOrgUnit(orgUnit, operatorID); err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := dto.OrgUnitResponse{
		ID:       orgUnit.ID,
		Name:     orgUnit.Name,
		ParentID: orgUnit.ParentID,
		Path:     orgUnit.Path,
		Level:    orgUnit.Level,
	}

	result.SuccessResponse(c, "组织节点更新成功", &response)
}

// DeleteOrgUnit 删除组织节点
// @Summary 删除组织节点
// @Description 删除组织节点
// @Tags 组织
// @Accept json
// @Produce json
// @Param id path int true "组织ID"
// @Success 200 {object} result.ResponseResult[string] "删除成功"
// @Failure 400 {object} result.ResponseResult[string] "请求参数错误"
// @Failure 500 {object} result.ResponseResult[string] "服务器内部错误"
// @Router /org-units/:id [delete]
func (h *OrgUnitHandler) DeleteOrgUnit(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	operatorID := c.GetUint("user_id")
	if err := h.orgService.DeleteOrgUnit(uint(id), operatorID); err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.SimpleSuccessResponse(c, "组织节点删除成功")
}

// GetOrgUnit 获取组织节点详情
// @Summary 获取组织节点详情
// @Description 根据ID获取组织节点详情
// @Tags 组织
// @Produce json
// @Param id path int true "组织ID"
// @Success 200 {object} result.ResponseResult[dto.OrgUnitResponse] "获取成功"
// @Failure 404 {object} result.ResponseResult[string] "组织不存在"
// @Router /org-units/:id [get]
func (h *OrgUnitHandler) GetOrgUnit(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	orgUnit, err := h.orgService.GetOrgUnitByID(uint(id))
	if err != nil {
		result.ErrorResponse(c, http.StatusNotFound, "组织节点不存在")
		return
	}

	response := dto.OrgUnitResponse{
		ID:       orgUnit.ID,
		Name:     orgUnit.Name,
		ParentID: orgUnit.ParentID,
		Path:     orgUnit.Path,
		Level:    orgUnit.Level,
	}

	result.SuccessResponse(c, "获取组织节点成功", &response)
}

// ListOrgUnits 获取组织节点列表
// @Summary 获取组织节点列表
// @Description 获取所有组织节点
// @Tags 组织
// @Produce json
// @Success 200 {object} result.ResponseResult[[]dto.OrgUnitResponse] "获取成功"
// @Router /org-units [get]
func (h *OrgUnitHandler) ListOrgUnits(c *gin.Context) {
	orgs, err := h.orgService.GetAllOrgUnits()
	if err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	responses := make([]dto.OrgUnitResponse, len(orgs))
	for i, org := range orgs {
		responses[i] = dto.OrgUnitResponse{
			ID:       org.ID,
			Name:     org.Name,
			ParentID: org.ParentID,
			Path:     org.Path,
			Level:    org.Level,
		}
	}

	result.SuccessResponse(c, "获取组织节点列表成功", &responses)
}

// GetOrgTree 获取组织树
// @Summary 获取组织树
// @Description 获取组织结构树
// @Tags 组织
// @Produce json
// @Success 200 {object} result.ResponseResult[[]dto.OrgUnitResponse] "获取成功"
// @Router /org-units/tree [get]
func (h *OrgUnitHandler) GetOrgTree(c *gin.Context) {
	orgs, err := h.orgService.GetOrgTree()
	if err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	responses := make([]dto.OrgUnitResponse, len(orgs))
	for i, org := range orgs {
		responses[i] = dto.OrgUnitResponse{
			ID:       org.ID,
			Name:     org.Name,
			ParentID: org.ParentID,
			Path:     org.Path,
			Level:    org.Level,
		}
	}

	result.SuccessResponse(c, "获取组织树成功", &responses)
}
