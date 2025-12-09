package entity

import "github.com/lyj404/gin-api-template/global"

type User struct {
	global.G_MODEL
	Name     string `gorm:"type:varchar(15);"`
	Email    string `gorm:"type:varchar(50);"`
	PassWord string `gorm:"type:varchar(100);column:password;"`
}

func (User) TableName() string {
	return "user"
}
