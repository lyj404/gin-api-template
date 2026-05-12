package entity

import "github.com/lyj404/gin-api-template/global"

// Role 角色实体
type Role struct {
	global.G_MODEL
	Name          string         `gorm:"type:varchar(100);unique;not null" json:"name"`        // 角色名称
	Description   string         `gorm:"type:varchar(255)" json:"description"`                // 角色描述
	IsSystem      bool           `gorm:"default:false" json:"is_system"`                     // 是否是系统内置角色
	RoleResources []RoleResource `gorm:"-" json:"role_resources,omitempty" binding:"-"`       // 角色资源（运行时填充）
	RoleOrgScopes []RoleOrgScope `gorm:"-" json:"role_org_scopes,omitempty" binding:"-"`      // 角色组织范围（运行时填充）
}
