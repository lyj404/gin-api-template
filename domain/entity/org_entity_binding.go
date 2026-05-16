package entity

import "github.com/lyj404/gin-api-template/global"

// OrgEntityBinding 组织实体绑定关系，用于数据权限控制
type OrgEntityBinding struct {
	global.G_MODEL
	OrgUnitID  uint64    `gorm:"not null;index" json:"org_unit_id"`             // 组织节点ID
	EntityType string  `gorm:"type:varchar(50);not null;index" json:"entity_type"` // 业务实体类型（如 device, order）
	EntityID   uint64    `gorm:"not null;index" json:"entity_id"`               // 业务实体ID
	OrgUnit    OrgUnit `gorm:"foreignKey:OrgUnitID" json:"-"`                 // 关联组织（不返回）
}
