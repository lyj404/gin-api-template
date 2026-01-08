package entity

import "github.com/lyj404/gin-api-template/global"

type Resource struct {
	global.G_MODEL
	Name        string `gorm:"unique;not null;size:100" json:"name"` // 资源名称/标识符
	Type        string `gorm:"not null;size:20" json:"type"`         // 资源类型: api/entity
	Pattern     string `gorm:"not null;index" json:"pattern"`        // API路径或实体模式（支持通配符）
	Method      string `gorm:"size:10" json:"method"`                // HTTP方法（仅API类型）
	Entity      string `gorm:"size:100" json:"entity"`               // 业务实体名称（仅Entity类型）
	Action      string `gorm:"size:20" json:"action"`                // 操作类型（仅Entity类型）
	Description string `gorm:"size:255" json:"description"`          // 资源描述
}
