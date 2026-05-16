package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lyj404/gin-api-template/domain/dto"
	"github.com/lyj404/gin-api-template/domain/entity"
	"github.com/lyj404/gin-api-template/domain/result"
	"github.com/lyj404/gin-api-template/domain/services"
)

type DictionaryHandler struct {
	dictService services.DictionaryService
}

func NewDictionaryHandler(dictService services.DictionaryService) *DictionaryHandler {
	return &DictionaryHandler{dictService: dictService}
}

// CreateDict 创建字典类型
func (h *DictionaryHandler) CreateDict(c *gin.Context) {
	var req dto.CreateDictRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	dict := &entity.SysDictionary{
		Name:   req.Name,
		Type:   req.Type,
		Status: req.Status,
		Desc:   req.Desc,
	}
	if dict.Status == 0 {
		dict.Status = 1
	}

	if err := h.dictService.CreateDict(c.Request.Context(), dict); err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.SuccessResponse(c, "创建成功", dict)
}

// UpdateDict 更新字典类型
func (h *DictionaryHandler) UpdateDict(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	dict, err := h.dictService.GetDictByID(c.Request.Context(), strconv.FormatUint(id, 10))
	if err != nil {
		result.ErrorResponse(c, http.StatusNotFound, "字典不存在")
		return
	}

	var req dto.UpdateDictRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if req.Name != "" {
		dict.Name = req.Name
	}
	if req.Type != "" {
		dict.Type = req.Type
	}
	if req.Status != 0 {
		dict.Status = req.Status
	}
	dict.Desc = req.Desc

	if err := h.dictService.UpdateDict(c.Request.Context(), &dict); err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.SimpleSuccessResponse(c, "更新成功")
}

// DeleteDict 删除字典类型
func (h *DictionaryHandler) DeleteDict(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := h.dictService.DeleteDict(c.Request.Context(), strconv.FormatUint(id, 10)); err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.SimpleSuccessResponse(c, "删除成功")
}

// GetDict 获取字典详情
func (h *DictionaryHandler) GetDict(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	dict, err := h.dictService.GetDictByID(c.Request.Context(), strconv.FormatUint(id, 10))
	if err != nil {
		result.ErrorResponse(c, http.StatusNotFound, "字典不存在")
		return
	}

	result.SuccessResponse(c, "获取成功", &dict)
}

// ListDict 字典列表
func (h *DictionaryHandler) ListDict(c *gin.Context) {
	name := c.Query("name")
	dictType := c.Query("type")

	dicts, err := h.dictService.ListDict(c.Request.Context(), name, dictType)
	if err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.SuccessResponse(c, "获取成功", &dicts)
}

// CreateDictDetail 创建字典详情
func (h *DictionaryHandler) CreateDictDetail(c *gin.Context) {
	var req dto.CreateDictDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	detail := &entity.SysDictionaryDetail{
		DictID: req.DictID,
		Label:  req.Label,
		Value:  req.Value,
		Sort:   req.Sort,
		Status: req.Status,
		Remark: req.Remark,
	}
	if detail.Status == 0 {
		detail.Status = 1
	}

	if err := h.dictService.CreateDictDetail(c.Request.Context(), detail); err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.SuccessResponse(c, "创建成功", detail)
}

// UpdateDictDetail 更新字典详情
func (h *DictionaryHandler) UpdateDictDetail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("detailId"), 10, 64)

	detail, err := h.dictService.GetDictDetailByID(c.Request.Context(), strconv.FormatUint(id, 10))
	if err != nil {
		result.ErrorResponse(c, http.StatusNotFound, "字典详情不存在")
		return
	}

	var req dto.UpdateDictDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if req.Label != "" {
		detail.Label = req.Label
	}
	if req.Value != "" {
		detail.Value = req.Value
	}
	if req.Sort != 0 {
		detail.Sort = req.Sort
	}
	if req.Status != 0 {
		detail.Status = req.Status
	}
	detail.Remark = req.Remark

	if err := h.dictService.UpdateDictDetail(c.Request.Context(), &detail); err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.SimpleSuccessResponse(c, "更新成功")
}

// DeleteDictDetail 删除字典详情
func (h *DictionaryHandler) DeleteDictDetail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("detailId"), 10, 64)

	if err := h.dictService.DeleteDictDetail(c.Request.Context(), strconv.FormatUint(id, 10)); err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.SimpleSuccessResponse(c, "删除成功")
}

// ListDictDetails 获取字典详情列表
func (h *DictionaryHandler) ListDictDetails(c *gin.Context) {
	dictID := c.Param("id")

	details, err := h.dictService.ListDictDetails(c.Request.Context(), dictID)
	if err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.SuccessResponse(c, "获取成功", &details)
}

// GetDictInfoByType 公共接口：根据类型获取字典详情
func (h *DictionaryHandler) GetDictInfoByType(c *gin.Context) {
	dictType := c.Param("type")

	details, err := h.dictService.GetDictInfoByType(c.Request.Context(), dictType)
	if err != nil {
		result.ErrorResponse(c, http.StatusNotFound, "字典类型不存在")
		return
	}

	result.SuccessResponse(c, "获取成功", &details)
}
