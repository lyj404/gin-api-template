package entity

import "github.com/lyj404/gin-api-template/global"

type UserRole struct {
	global.G_MODEL
	UserID    uint    `gorm:"not null;index" json:"user_id"`     // 用户ID
	RoleID    uint    `gorm:"not null;index" json:"role_id"`     // 角色ID
	OrgUnitID uint    `gorm:"not null;index" json:"org_unit_id"` // 角色生效的组织节点
	User      User    `gorm:"foreignKey:UserID"`
	Role      Role    `gorm:"foreignKey:RoleID"`
	OrgUnit   OrgUnit `gorm:"foreignKey:OrgUnitID"`
}
