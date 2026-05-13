package repositories

import (
	"github.com/lyj404/gin-api-template/domain/entity"
	"gorm.io/gorm"
)

// UserRepository 用户管理仓储接口（区别于 domain.UserRepo 主要用于登录场景）
type UserRepository interface {
	List(page, pageSize int, keyword string) ([]entity.User, int64, error)
	GetByID(id uint) (*entity.User, error)
	Create(tx *gorm.DB, user *entity.User) error
	Update(tx *gorm.DB, user *entity.User) error
	UpdatePassword(tx *gorm.DB, id uint, hashed string) error
	Delete(tx *gorm.DB, id uint) error
	GetRoleIDsByUserID(userID uint) ([]uint, error)
	ReplaceUserRoles(tx *gorm.DB, userID, orgUnitID uint, roleIDs []uint) error
	GetRoleNamesByUserID(userID uint) ([]string, error)
}
