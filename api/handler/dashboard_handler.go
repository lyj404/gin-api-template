package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lyj404/gin-api-template/domain/result"
	"github.com/lyj404/gin-api-template/domain/services"
)

// DashboardHandler 仪表盘处理器，依赖 DashboardService 进行权限过滤
type DashboardHandler struct {
	dashboardSvc services.DashboardService
}

func NewDashboardHandler(dashboardSvc services.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		dashboardSvc: dashboardSvc,
	}
}

// Stats 仪表盘统计
// @Summary 仪表盘汇总统计
// @Description 获取用户、角色、菜单、资源的总数统计（根据当前用户权限过滤）
// @Tags 仪表盘
// @Produce json
// @Success 200 {object} result.ResponseResult[services.DashboardStats] "获取成功"
// @Failure 500 {object} result.ResponseResult[string] "服务器内部错误"
// @Router /dashboard/stats [get]
func (h *DashboardHandler) Stats(c *gin.Context) {
	userID, _ := c.Get("user_id")

	stats, err := h.dashboardSvc.GetStats(userID.(uint64))
	if err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.SuccessResponse(c, "获取统计成功", stats)
}

// AuditTrend 最近 7 天审计日志趋势
// @Summary 审计日志 7 天趋势
// @Description 按日期分组返回最近 7 天的审计日志数量
// @Tags 仪表盘
// @Produce json
// @Success 200 {object} result.ResponseResult[[]services.AuditTrendItem] "获取成功"
// @Failure 500 {object} result.ResponseResult[string] "服务器内部错误"
// @Router /dashboard/audit-trend [get]
func (h *DashboardHandler) AuditTrend(c *gin.Context) {
	userID, _ := c.Get("user_id")
	items, err := h.dashboardSvc.GetAuditTrend(userID.(uint64))
	if err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.SuccessResponse(c, "获取审计趋势成功", &items)
}
