package service

import (
	"github.com/lyj404/gin-api-template/domain/entity"
	"gorm.io/gorm"
)

func getOperatorName(tx *gorm.DB, userID uint) string {
	var user entity.User
	if err := tx.First(&user, userID).Error; err != nil {
		return "未知"
	}
	return user.Name
}
