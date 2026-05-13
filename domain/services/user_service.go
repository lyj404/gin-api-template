package services

import (
	"github.com/lyj404/gin-api-template/domain/dto"
	"github.com/lyj404/gin-api-template/domain/entity"
)

// UserManagementService 用户管理服务接口（用户 CRUD + 角色分配）
type UserManagementService interface {
	List(page, pageSize int, keyword string) ([]entity.User, map[uint][]uint, map[uint][]string, int64, error)
	GetByID(id uint) (*entity.User, []uint, []string, error)
	Create(req *dto.CreateUserRequest, operatorID uint) (*entity.User, error)
	Update(id uint, req *dto.UpdateUserRequest, operatorID uint) (*entity.User, error)
	Delete(id uint, operatorID uint) error
}
