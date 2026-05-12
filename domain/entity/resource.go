package entity

import "github.com/lyj404/gin-api-template/global"

// Resource 资源实体，定义API路径或业务实体的权限
type Resource struct {
	global.G_MODEL
	Name        string `gorm:"type:varchar(100);unique;not null" json:"name"`        // 资源名称/标识符
	Type        string `gorm:"type:varchar(20);not null" json:"type"`               // 资源类型: api/entity
	Pattern     string `gorm:"type:varchar(255);not null;index" json:"pattern"`      // API路径或实体模式（支持通配符）
	Method      string `gorm:"type:varchar(10)" json:"method"`                       // HTTP方法（仅API类型）
	Entity      string `gorm:"type:varchar(100)" json:"entity"`                      // 业务实体名称（仅Entity类型）
	Action      string `gorm:"type:varchar(50)" json:"action"`                        // 操作类型（仅Entity类型）
	Description string `gorm:"type:varchar(255)" json:"description"`                 // 资源描述
}
