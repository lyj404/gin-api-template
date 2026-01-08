package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lyj404/gin-api-template/domain/entity"
	"github.com/lyj404/gin-api-template/domain/result"
	domainservices "github.com/lyj404/gin-api-template/domain/services"
)

type AuditLogHandler struct {
	auditLogService domainservices.AuditLogService
}

func NewAuditLogHandler(auditLogService domainservices.AuditLogService) *AuditLogHandler {
	return &AuditLogHandler{
		auditLogService: auditLogService,
	}
}

// GetAuditLogsByOperator 按操作者查询审计日志
// @Summary 按操作者查询审计日志
// @Description 根据操作者ID分页查询审计日志
// @Tags 审计
// @Produce json
// @Param operator_id query int true "操作者ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} result.ResponseResult[PageResponse] "查询成功"
// @Failure 500 {object} result.ResponseResult[string] "服务器内部错误"
// @Router /audit-logs [get]
func (h *AuditLogHandler) GetAuditLogsByOperator(c *gin.Context) {
	operatorID, _ := strconv.Atoi(c.Query("operator_id"))
	if operatorID == 0 {
		result.ErrorResponse(c, http.StatusBadRequest, "缺少操作者ID")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	logs, total, err := h.auditLogService.GetAuditLogsByOperator(uint(operatorID), page, pageSize)
	if err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := PageResponse{
		Data:     logs,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}

	result.SuccessResponse(c, "查询审计日志成功", &response)
}

// GetAuditLogsByTarget 按目标查询审计日志
// @Summary 按目标查询审计日志
// @Description 根据目标类型和ID分页查询审计日志
// @Tags 审计
// @Produce json
// @Param target_type query string true "目标类型"
// @Param target_id query int true "目标ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} result.ResponseResult[PageResponse] "查询成功"
// @Failure 500 {object} result.ResponseResult[string] "服务器内部错误"
// @Router /audit-logs/target [get]
func (h *AuditLogHandler) GetAuditLogsByTarget(c *gin.Context) {
	targetType := c.Query("target_type")
	if targetType == "" {
		result.ErrorResponse(c, http.StatusBadRequest, "缺少目标类型")
		return
	}

	targetID, _ := strconv.Atoi(c.Query("target_id"))
	if targetID == 0 {
		result.ErrorResponse(c, http.StatusBadRequest, "缺少目标ID")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	logs, total, err := h.auditLogService.GetAuditLogsByTarget(targetType, uint(targetID), page, pageSize)
	if err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := PageResponse{
		Data:     logs,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}

	result.SuccessResponse(c, "查询审计日志成功", &response)
}

// GetAuditLogsByTimeRange 按时间范围查询审计日志
// @Summary 按时间范围查询审计日志
// @Description 根据时间范围分页查询审计日志
// @Tags 审计
// @Produce json
// @Param start_time query string false "开始时间"
// @Param end_time query string false "结束时间"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} result.ResponseResult[PageResponse] "查询成功"
// @Failure 500 {object} result.ResponseResult[string] "服务器内部错误"
// @Router /audit-logs/time [get]
func (h *AuditLogHandler) GetAuditLogsByTimeRange(c *gin.Context) {
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	logs, total, err := h.auditLogService.GetAuditLogsByTimeRange(startTime, endTime, page, pageSize)
	if err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := PageResponse{
		Data:     logs,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}

	result.SuccessResponse(c, "查询审计日志成功", &response)
}

type PageResponse struct {
	Data     []entity.AuditLog `json:"data"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
}
