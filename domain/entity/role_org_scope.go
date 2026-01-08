package entity

import "github.com/lyj404/gin-api-template/global"

type RoleOrgScope struct {
	global.G_MODEL
	RoleID             uint    `gorm:"not null;index" json:"role_id"`                     // 角色ID
	OrgUnitID          uint    `gorm:"not null;index" json:"org_unit_id"`                 // 组织节点ID
	IncludeDescendants bool    `gorm:"not null;default:false" json:"include_descendants"` // 是否包含所有子节点
	Role               Role    `gorm:"foreignKey:RoleID"`
	OrgUnit            OrgUnit `gorm:"foreignKey:OrgUnitID"`
}
