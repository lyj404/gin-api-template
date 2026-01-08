package entity

import "github.com/lyj404/gin-api-template/global"

type OrgUnit struct {
	global.G_MODEL
	Name     string `gorm:"not null;size:100" json:"name"`       // 组织名称
	ParentID *uint  `gorm:"index" json:"parent_id"`              // 父节点ID
	Path     string `gorm:"not null;size:500;index" json:"path"` // 组织路径（如 /1/3/5）
	Level    int    `gorm:"not null;default:0" json:"level"`     // 层级深度（root=0）
}
