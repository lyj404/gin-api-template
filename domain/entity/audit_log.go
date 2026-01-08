package entity

import "github.com/lyj404/gin-api-template/global"

type AuditLog struct {
	global.G_MODEL
	OperatorID   uint   `gorm:"not null;index" json:"operator_id"`         // 操作者ID
	OperatorName string `gorm:"not null;size:100" json:"operator_name"`    // 操作者姓名
	Action       string `gorm:"not null;size:50;index" json:"action"`      // 操作类型（create, update, delete, assign）
	TargetType   string `gorm:"not null;size:50;index" json:"target_type"` // 目标类型（role, resource, org, user_role）
	TargetID     uint   `gorm:"not null;index" json:"target_id"`           // 目标ID
	BeforeData   string `gorm:"type:text" json:"before_data"`              // 变更前数据（JSON）
	AfterData    string `gorm:"type:text" json:"after_data"`               // 变更后数据（JSON）
	Description  string `gorm:"type:text" json:"description"`              // 操作描述
}
