package services

import (
	"github.com/lyj404/gin-api-template/domain/dto"
	"github.com/lyj404/gin-api-template/domain/entity"
)

// UserManagementService 用户管理服务接口（用户 CRUD + 角色分配）
type UserManagementService interface {
	List(page, pageSize int, keyword string, userID uint64) ([]entity.User, map[uint64][]uint64, map[uint64][]string, int64, error)
	GetByID(id uint64, operatorID uint64) (*entity.User, []uint64, []string, error)
	Create(req *dto.CreateUserRequest, operatorID uint64) (*entity.User, error)
	Update(id uint64, req *dto.UpdateUserRequest, operatorID uint64) (*entity.User, error)
	Delete(id uint64, operatorID uint64) error
}
