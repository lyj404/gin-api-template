package entity

import "github.com/lyj404/gin-api-template/global"

type RoleMenu struct {
	global.G_MODEL
	RoleID uint64  `gorm:"not null;index" json:"role_id"`
	MenuID uint64  `gorm:"not null;index" json:"menu_id"`
	Role   *Role `gorm:"foreignKey:RoleID" json:"-"`
	Menu   *Menu `gorm:"foreignKey:MenuID" json:"-"`
}
