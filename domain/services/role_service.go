package services

import (
	"github.com/lyj404/gin-api-template/domain/entity"
)

type RoleService interface {
	CreateRole(role *entity.Role, operatorID uint) error
	UpdateRole(role *entity.Role, operatorID uint) error
	DeleteRole(id uint, operatorID uint) error
	GetRoleByID(id uint) (*entity.Role, error)
	GetAllRoles() ([]entity.Role, error)
	BindResource(roleID, resourceID uint, isWrite bool, operatorID uint) error
	UnbindResource(roleID, resourceID uint, operatorID uint) error
	BindOrgScope(roleID, orgUnitID uint, includeDescendants bool, operatorID uint) error
	UnbindOrgScope(roleID, orgUnitID uint, operatorID uint) error
	AssignRoleToUser(userID, roleID, orgUnitID uint, operatorID uint) error
	RevokeRoleFromUser(userID, roleID, orgUnitID uint, operatorID uint) error
}
