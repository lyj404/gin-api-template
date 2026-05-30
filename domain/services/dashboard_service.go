package services

// DashboardStats 仪表盘统计，用户只能看到其权限范围内的数据
type DashboardStats struct {
	UserCount     int64 `json:"user_count"`
	RoleCount     int64 `json:"role_count"`
	MenuCount     int64 `json:"menu_count"`
	ResourceCount int64 `json:"resource_count"`
}

// AuditTrendItem 审计日志趋势项
type AuditTrendItem struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

// DashboardService 仪表盘服务接口，所有统计根据用户权限/组织范围过滤
type DashboardService interface {
	GetStats(userID uint64) (*DashboardStats, error)
	GetAuditTrend(userID uint64) ([]AuditTrendItem, error)
}
