package entity

import "github.com/lyj404/gin-api-template/global"

type OrgEntityBinding struct {
	global.G_MODEL
	OrgUnitID  uint    `gorm:"not null;index" json:"org_unit_id"`         // 组织节点ID
	EntityType string  `gorm:"not null;size:50;index" json:"entity_type"` // 业务实体类型（如 device, order）
	EntityID   uint    `gorm:"not null;index" json:"entity_id"`           // 业务实体ID
	OrgUnit    OrgUnit `gorm:"foreignKey:OrgUnitID"`
}
