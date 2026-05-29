package service

import (
	"fmt"
	"time"

	"github.com/lyj404/gin-api-template/domain/entity"
	"github.com/lyj404/gin-api-template/domain/services"
	"github.com/lyj404/gin-api-template/global"
)

type dashboardServiceImpl struct {
	permSvc services.PermissionService
}

func NewDashboardService(permSvc services.PermissionService) services.DashboardService {
	return &dashboardServiceImpl{
		permSvc: permSvc,
	}
}

func (s *dashboardServiceImpl) GetStats(userID uint64) (*services.DashboardStats, error) {
	isSuper, err := s.permSvc.HasSystemRole(userID)
	if err != nil {
		return nil, fmt.Errorf("检查系统角色失败: %w", err)
	}

	var stats services.DashboardStats

	if isSuper {
		if err := global.G_DB.Model(&entity.UserRole{}).
			Distinct("user_id").
			Count(&stats.UserCount).Error; err != nil {
			return nil, err
		}
		if err := global.G_DB.Model(&entity.RoleOrgScope{}).
			Distinct("role_id").
			Count(&stats.RoleCount).Error; err != nil {
			return nil, err
		}
	} else {
		orgScope, err := s.permSvc.GetUserOrgScope(userID)
		if err != nil {
			return nil, fmt.Errorf("获取组织范围失败: %w", err)
		}
		orgIDs := CollectOrgIDs(orgScope)

		if err := global.G_DB.Model(&entity.UserRole{}).
			Where("org_unit_id IN ?", orgIDs).
			Distinct("user_id").
			Count(&stats.UserCount).Error; err != nil {
			return nil, err
		}
		if err := global.G_DB.Model(&entity.RoleOrgScope{}).
			Where("org_unit_id IN ?", orgIDs).
			Distinct("role_id").
			Count(&stats.RoleCount).Error; err != nil {
			return nil, err
		}
	}

	menus, err := s.permSvc.GetUserMenus(userID)
	if err != nil {
		return nil, err
	}
	stats.MenuCount = countMenuNodes(menus)

	perms, err := s.permSvc.GetUserPermissions(userID)
	if err != nil {
		return nil, err
	}
	stats.ResourceCount = int64(len(perms))

	return &stats, nil
}

func (s *dashboardServiceImpl) GetAuditTrend() ([]services.AuditTrendItem, error) {
	type row struct {
		LogDate string `gorm:"column:log_date"`
		Count   int64
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
		Select(dateExpr+" as log_date, COUNT(*) as count").
		Where("created_at >= ?", cutoff).
		Group("log_date").
		Order("log_date ASC").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	out := make([]services.AuditTrendItem, 0, len(rows))
	for _, r := range rows {
		out = append(out, services.AuditTrendItem{Date: r.LogDate, Count: r.Count})
	}
	return out, nil
}

func CollectOrgIDs(orgScope []services.OrgScopeInfo) []uint64 {
	idSet := make(map[uint64]struct{})
	for _, scope := range orgScope {
		if scope.IncludeDescendants {
			var ids []uint64
			global.G_DB.Model(&entity.OrgUnit{}).
				Where("path LIKE ? OR id = ?", scope.Path+"/%", scope.OrgUnitID).
				Pluck("id", &ids)
			for _, id := range ids {
				idSet[id] = struct{}{}
			}
		} else {
			idSet[scope.OrgUnitID] = struct{}{}
		}
	}
	ids := make([]uint64, 0, len(idSet))
	for id := range idSet {
		ids = append(ids, id)
	}
	return ids
}

func countMenuNodes(menus []services.MenuTreeNode) int64 {
	var count int64
	for _, m := range menus {
		count++
		if len(m.Children) > 0 {
			count += countMenuNodes(m.Children)
		}
	}
	return count
}
