package entity

import "github.com/lyj404/gin-api-template/global"

type Role struct {
	global.G_MODEL
	Name          string         `gorm:"unique;not null;size:100" json:"name"` // 角色名称
	Description   string         `gorm:"size:255" json:"description"`          // 角色描述
	IsSystem      bool           `gorm:"default:false" json:"is_system"`       // 是否是系统内置角色
	RoleResources []RoleResource `gorm:"foreignKey:RoleID"`
	RoleOrgScopes []RoleOrgScope `gorm:"foreignKey:RoleID"`
}
