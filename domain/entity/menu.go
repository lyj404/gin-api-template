package entity

import "github.com/lyj404/gin-api-template/global"

type Menu struct {
	global.G_MODEL
	Name       string     `gorm:"type:varchar(100);not null" json:"name"`
	ParentID   *uint      `gorm:"index" json:"parent_id"`
	Path       string     `gorm:"type:varchar(255)" json:"path"`
	Icon       string     `gorm:"type:varchar(100)" json:"icon"`
	OrderNum   int        `gorm:"default:0" json:"order_num"`
	IsVisible  bool       `gorm:"default:true" json:"is_visible"`
	Status     string     `gorm:"type:varchar(20);default:enabled" json:"status"`
	Children   []Menu     `gorm:"-" json:"children,omitempty" binding:"-"`
	Resources  []Resource `gorm:"many2many:menu_resource;" json:"resources,omitempty"`
}
