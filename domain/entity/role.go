package entity

import "github.com/lyj404/gin-api-template/global"

type Role struct {
	global.G_MODEL
	Name          string         `gorm:"type:varchar(100);unique;not null" json:"name"`
	Description   string         `gorm:"type:varchar(255)" json:"description"`
	IsSystem      bool           `gorm:"default:false" json:"is_system"`
	RoleResources []RoleResource `gorm:"foreignKey:RoleID" json:"role_resources,omitempty" binding:"-"`
	RoleOrgScopes []RoleOrgScope `gorm:"foreignKey:RoleID" json:"role_org_scopes,omitempty" binding:"-"`
	RoleMenus     []RoleMenu     `gorm:"foreignKey:RoleID" json:"role_menus,omitempty" binding:"-"`
}
