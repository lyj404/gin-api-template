package entity

import "github.com/lyj404/gin-api-template/global"

// Menu 菜单实体，关联资源以控制菜单显示权限
type Menu struct {
	global.G_MODEL
	Name        string  `gorm:"not null;size:100" json:"name"`                        // 菜单名称
	ParentID    *uint   `gorm:"index" json:"parent_id"`                              // 父菜单ID，nil表示根菜单
	Path        string  `gorm:"size:255" json:"path"`                                 // 前端路由路径
	Component   string  `gorm:"size:255" json:"component"`                            // 前端组件路径
	Icon        string  `gorm:"size:100" json:"icon"`                                 // 菜单图标
	OrderNum    int     `gorm:"default:0" json:"order_num"`                          // 排序序号
	ResourceID  uint    `gorm:"not null;index" json:"resource_id"`                    // 关联资源ID，通过资源权限控制菜单显示
	IsVisible   bool    `gorm:"default:true" json:"is_visible"`                      // 是否显示
	Status      string  `gorm:"default:enabled;size:20" json:"status"`                // 状态：enabled/disabled
	Children    []Menu  `gorm:"-" json:"children,omitempty" binding:"-"`              // 子菜单（运行时填充）
	Resource    *Resource `gorm:"foreignKey:ResourceID" json:"resource,omitempty"`   // 关联资源信息
}