package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lyj404/gin-api-template/domain/dto"
	"github.com/lyj404/gin-api-template/domain/result"
	domainservices "github.com/lyj404/gin-api-template/domain/services"
)

// AuditLogHandler 审计日志处理器，处理审计日志相关的HTTP请求
type AuditLogHandler struct {
	auditLogService domainservices.AuditLogService
}

// NewAuditLogHandler 创建审计日志处理器实例
func NewAuditLogHandler(auditLogService domainservices.AuditLogService) *AuditLogHandler {
	return &AuditLogHandler{
		auditLogService: auditLogService,
	}
}

// GetAuditLogsByOperator 按操作者查询审计日志，不传操作者ID时返回全部
// @Summary 按操作者查询审计日志
// @Description 根据操作者ID分页查询审计日志，不传操作者ID时返回全部日志
// @Tags 审计
// @Produce json
// @Param operator_id query int false "操作者ID，不传时返回全部"
// @Param page query int false "页码，默认1"
// @Param page_size query int false "每页数量，默认10，最大100"
// @Success 200 {object} result.ResponseResult[dto.PaginationResponse] "查询成功"
// @Failure 400 {object} result.ResponseResult[string] "参数错误"
// @Failure 500 {object} result.ResponseResult[string] "服务器内部错误"
// @Router /audit-logs [get]
func (h *AuditLogHandler) GetAuditLogsByOperator(c *gin.Context) {
	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		result.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	req.SetDefaults()

	operatorID, _ := strconv.ParseUint(c.Query("operator_id"), 10, 64)
	if operatorID != 0 {
		logs, total, err := h.auditLogService.GetAuditLogsByOperator(operatorID, req.Page, req.PageSize)
		if err != nil {
			result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		result.SuccessResponse(c, "查询审计日志成功", dto.NewPaginationResponse(
			req.Page,
			req.PageSize,
			total,
			logs,
		))
		return
	}

	// 不传操作者ID时返回全部日志
	logs, total, err := h.auditLogService.GetAuditLogsByTimeRange("", "", req.Page, req.PageSize)
	if err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	result.SuccessResponse(c, "查询审计日志成功", dto.NewPaginationResponse(
		req.Page,
		req.PageSize,
		total,
		logs,
	))
}

// GetAuditLogsByTarget 按目标查询审计日志
// @Summary 按目标查询审计日志
// @Description 根据目标类型和ID分页查询审计日志
// @Tags 审计
// @Produce json
// @Param target_type query string true "目标类型"
// @Param target_id query int true "目标ID"
// @Param page query int false "页码，默认1"
// @Param page_size query int false "每页数量，默认10，最大100"
// @Success 200 {object} result.ResponseResult[dto.PaginationResponse] "查询成功"
// @Failure 400 {object} result.ResponseResult[string] "缺少参数"
// @Failure 500 {object} result.ResponseResult[string] "服务器内部错误"
// @Router /audit-logs/target [get]
func (h *AuditLogHandler) GetAuditLogsByTarget(c *gin.Context) {
	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		result.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	req.SetDefaults()

	targetType := c.Query("target_type")
	if targetType == "" {
		result.ErrorResponse(c, http.StatusBadRequest, "缺少目标类型")
		return
	}

	targetID, _ := strconv.ParseUint(c.Query("target_id"), 10, 64)
	if targetID == 0 {
		result.ErrorResponse(c, http.StatusBadRequest, "缺少目标ID")
		return
	}

	logs, total, err := h.auditLogService.GetAuditLogsByTarget(targetType, targetID, req.Page, req.PageSize)
	if err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.SuccessResponse(c, "查询审计日志成功", dto.NewPaginationResponse(
		req.Page,
		req.PageSize,
		total,
		logs,
	))
}

// GetAuditLogsByTimeRange 按时间范围查询审计日志
// @Summary 按时间范围查询审计日志
// @Description 根据时间范围分页查询审计日志
// @Tags 审计
// @Produce json
// @Param start_time query string false "开始时间"
// @Param end_time query string false "结束时间"
// @Param page query int false "页码，默认1"
// @Param page_size query int false "每页数量，默认10，最大100"
// @Success 200 {object} result.ResponseResult[dto.PaginationResponse] "查询成功"
// @Failure 500 {object} result.ResponseResult[string] "服务器内部错误"
// @Router /audit-logs/time [get]
func (h *AuditLogHandler) GetAuditLogsByTimeRange(c *gin.Context) {
	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		result.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	req.SetDefaults()

	startTime := c.Query("start_time")
	endTime := c.Query("end_time")

	logs, total, err := h.auditLogService.GetAuditLogsByTimeRange(startTime, endTime, req.Page, req.PageSize)
	if err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.SuccessResponse(c, "查询审计日志成功", dto.NewPaginationResponse(
		req.Page,
		req.PageSize,
		total,
		logs,
	))
}
