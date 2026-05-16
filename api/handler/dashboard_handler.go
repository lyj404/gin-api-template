package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lyj404/gin-api-template/domain/entity"
	"github.com/lyj404/gin-api-template/domain/result"
	"github.com/lyj404/gin-api-template/global"
)

type DashboardHandler struct{}

func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{}
}

type dashboardStats struct {
	UserCount     int64 `json:"user_count"`
	RoleCount     int64 `json:"role_count"`
	MenuCount     int64 `json:"menu_count"`
	ResourceCount int64 `json:"resource_count"`
}

type auditTrendItem struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

// Stats 仪表盘统计
// @Summary 仪表盘汇总统计
// @Description 获取用户、角色、菜单、资源的总数统计
// @Tags 仪表盘
// @Produce json
// @Success 200 {object} result.ResponseResult[dashboardStats] "获取成功"
// @Failure 500 {object} result.ResponseResult[string] "服务器内部错误"
// @Router /dashboard/stats [get]
func (h *DashboardHandler) Stats(c *gin.Context) {
	var stats dashboardStats

	if err := global.G_DB.Model(&entity.User{}).Count(&stats.UserCount).Error; err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if err := global.G_DB.Model(&entity.Role{}).Count(&stats.RoleCount).Error; err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if err := global.G_DB.Model(&entity.Menu{}).Count(&stats.MenuCount).Error; err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if err := global.G_DB.Model(&entity.Resource{}).Count(&stats.ResourceCount).Error; err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.SuccessResponse(c, "获取统计成功", &stats)
}

// AuditTrend 最近 7 天审计日志趋势
// @Summary 审计日志 7 天趋势
// @Description 按日期分组返回最近 7 天的审计日志数量
// @Tags 仪表盘
// @Produce json
// @Success 200 {object} result.ResponseResult[[]auditTrendItem] "获取成功"
// @Failure 500 {object} result.ResponseResult[string] "服务器内部错误"
// @Router /dashboard/audit-trend [get]
func (h *DashboardHandler) AuditTrend(c *gin.Context) {
	type row struct {
		Date  string
		Count int64
	}
	var rows []row
	dialect := global.G_DB.Dialector.Name()
	var dateExpr string
	if dialect == "postgres" {
		dateExpr = "TO_CHAR(created_at, 'YYYY-MM-DD')"
	} else {
		dateExpr = "DATE_FORMAT(created_at, '%Y-%m-%d')"
	}

	cutoff := time.Now().AddDate(0, 0, -6).Truncate(24 * time.Hour)

	if err := global.G_DB.Model(&entity.AuditLog{}).
		Select(dateExpr+" as date, COUNT(*) as count").
		Where("created_at >= ?", cutoff).
		Group("date").
		Order("date ASC").
		Scan(&rows).Error; err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	out := make([]auditTrendItem, 0, len(rows))
	for _, r := range rows {
		out = append(out, auditTrendItem{Date: r.Date, Count: r.Count})
	}
	result.SuccessResponse(c, "获取审计趋势成功", &out)
}
