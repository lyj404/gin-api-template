package repositories

import (
	"github.com/lyj404/gin-api-template/domain/entity"
	"gorm.io/gorm"
)

// UserRepository 用户管理仓储接口（区别于 domain.UserRepo 主要用于登录场景）
type UserRepository interface {
	List(page, pageSize int, keyword string) ([]entity.User, int64, error)
	GetByID(id uint64) (*entity.User, error)
	Create(tx *gorm.DB, user *entity.User) error
	Update(tx *gorm.DB, user *entity.User) error
	UpdatePassword(tx *gorm.DB, id uint64, hashed string) error
	Delete(tx *gorm.DB, id uint64) error
	GetRoleIDsByUserID(userID uint64) ([]uint64, error)
	ReplaceUserRoles(tx *gorm.DB, userID, orgUnitID uint64, roleIDs []uint64) error
	GetRoleNamesByUserID(userID uint64) ([]string, error)
	HasSystemRole(userID uint64) (bool, error)
}
