package entity

import "github.com/lyj404/gin-api-template/global"

type MenuResource struct {
	global.G_MODEL
	MenuID     uint64      `gorm:"not null;index" json:"menu_id"`
	ResourceID uint64      `gorm:"not null;index" json:"resource_id"`
	Menu       *Menu     `gorm:"foreignKey:MenuID" json:"-"`
	Resource   *Resource `gorm:"foreignKey:ResourceID" json:"-"`
}
