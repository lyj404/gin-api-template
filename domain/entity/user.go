package entity

import "github.com/lyj404/gin-api-template/global"

// User 用户实体
type User struct {
	global.G_MODEL
	Name     string     `gorm:"type:varchar(50)" json:"name"`               // 用户姓名
	Email    string     `gorm:"type:varchar(255);unique;not null" json:"email"` // 用户邮箱
	PassWord string     `gorm:"type:varchar(255);column:password" json:"-"`  // 用户密码（不返回）
	Roles    []UserRole `gorm:"foreignKey:UserID" json:"-"`                  // 用户角色（不返回）
}

// TableName 指定表名为 user
func (User) TableName() string {
	return "user"
}
