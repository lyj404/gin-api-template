package entity

import "github.com/lyj404/gin-api-template/global"

// AuditLog 审计日志实体，记录权限变更操作
type AuditLog struct {
	global.G_MODEL
	OperatorID   uint64   `gorm:"not null;index" json:"operator_id"`          // 操作者ID
	OperatorName string `gorm:"type:varchar(100);not null" json:"operator_name"` // 操作者姓名
	Action       string `gorm:"type:varchar(50);not null;index" json:"action"`     // 操作类型（create, update, delete, assign）
	TargetType   string `gorm:"type:varchar(50);not null;index" json:"target_type"` // 目标类型（role, resource, org, user_role）
	TargetID     uint64   `gorm:"not null;index" json:"target_id"`            // 目标ID
	BeforeData   string `gorm:"type:text" json:"before_data"`               // 变更前数据（JSON）
	AfterData    string `gorm:"type:text" json:"after_data"`                // 变更后数据（JSON）
	Description  string `gorm:"type:text" json:"description"`               // 操作描述
}
