package entity

import "github.com/lyj404/gin-api-template/global"

// RoleResource 角色资源关联实体
type RoleResource struct {
	global.G_MODEL
	RoleID     uint64     `gorm:"not null;index" json:"role_id"`          // 角色ID
	ResourceID uint64     `gorm:"not null;index" json:"resource_id"`       // 资源ID
	IsRead     bool     `gorm:"not null;default:true" json:"is_read"`    // 是否有读权限（默认true）
	IsWrite    bool     `gorm:"not null;default:false" json:"is_write"`  // 是否有写权限
	Role       Role     `gorm:"foreignKey:RoleID" json:"-"`              // 关联角色（不返回）
	Resource   Resource `gorm:"foreignKey:ResourceID" json:"-"`           // 关联资源（不返回）
}
