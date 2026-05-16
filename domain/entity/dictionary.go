package entity

import "github.com/lyj404/gin-api-template/global"

// SysDictionary 字典类型实体
type SysDictionary struct {
	global.G_MODEL
	Name    string                `gorm:"column:name;type:varchar(100);not null" json:"name"`   // 字典名称
	Type    string                `gorm:"column:type;type:varchar(100);unique;not null" json:"type"` // 字典类型标识
	Status  int                   `gorm:"column:status;default:1" json:"status"`                 // 状态 (1:启用, 2:禁用)
	Desc    string                `gorm:"column:desc;type:varchar(255)" json:"desc"`              // 描述
	Details []SysDictionaryDetail `gorm:"foreignKey:DictID" json:"details"`                      // 字典详情
}

func (SysDictionary) TableName() string {
	return "sys_dictionary"
}

// SysDictionaryDetail 字典详情实体
type SysDictionaryDetail struct {
	global.G_MODEL
	DictID uint64 `gorm:"column:dict_id;not null" json:"dict_id"`              // 字典类型ID
	Label  string `gorm:"column:label;type:varchar(100);not null" json:"label"`             // 字典标签
	Value  string `gorm:"column:value;type:varchar(100);not null" json:"value"`             // 字典数值
	Sort   int    `gorm:"column:sort;default:0" json:"sort"`                                // 排序
	Status int    `gorm:"column:status;default:1" json:"status"`                            // 状态 (1:启用, 2:禁用)
	Remark string `gorm:"column:remark;type:varchar(255)" json:"remark"`                    // 备注
}

func (SysDictionaryDetail) TableName() string {
	return "sys_dictionary_detail"
}
