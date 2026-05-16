package entity

import "github.com/lyj404/gin-api-template/global"

// RoleOrgScope 角色组织范围关联实体
type RoleOrgScope struct {
	global.G_MODEL
	RoleID             uint64    `gorm:"not null;index" json:"role_id"`                      // 角色ID
	OrgUnitID          uint64    `gorm:"not null;index" json:"org_unit_id"`                  // 组织节点ID
	IncludeDescendants bool    `gorm:"not null;default:false" json:"include_descendants"`   // 是否包含所有子节点
	Role               Role    `gorm:"foreignKey:RoleID" json:"-"`                          // 关联角色（不返回）
	OrgUnit            OrgUnit `gorm:"foreignKey:OrgUnitID" json:"-"`                        // 关联组织（不返回）
}
