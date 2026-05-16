package entity

import "github.com/lyj404/gin-api-template/global"

// OrgUnit 组织单元实体，树形结构
type OrgUnit struct {
	global.G_MODEL
	Name     string `gorm:"type:varchar(100);not null" json:"name"`              // 组织名称
	ParentID *uint64  `gorm:"index" json:"parent_id"`                             // 父节点ID
	Path     string `gorm:"type:varchar(1000);not null;index" json:"path"`      // 组织路径（如 /1/3/5）
	Level    int    `gorm:"not null;default:0" json:"level"`                    // 层级深度（root=0）
}
