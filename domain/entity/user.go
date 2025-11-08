package entity

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"type:varchar(15);"`
	Email    string `gorm:"type:varchar(50);"`
	PassWord string `gorm:"type:varchar(100);column:password;"`
}

func (User) TableName() string {
	return "user"
}
