package entity

import "github.com/lyj404/gin-api-template/global"

// UserRole 用户角色关联实体
type UserRole struct {
	global.G_MODEL
	UserID    uint    `gorm:"not null;index" json:"user_id"`     // 用户ID
	RoleID    uint    `gorm:"not null;index" json:"role_id"`     // 角色ID
	OrgUnitID uint    `gorm:"not null;index" json:"org_unit_id"` // 角色生效的组织节点
	User      User    `gorm:"foreignKey:UserID" json:"-"`        // 关联用户（不返回）
	Role      Role    `gorm:"foreignKey:RoleID" json:"-"`       // 关联角色（不返回）
	OrgUnit   OrgUnit `gorm:"foreignKey:OrgUnitID" json:"-"`     // 关联组织（不返回）
}
