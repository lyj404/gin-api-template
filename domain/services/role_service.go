package services

import (
	"github.com/lyj404/gin-api-template/domain/entity"
)

type RoleService interface {
	CreateRole(role *entity.Role, operatorID uint64) error
	UpdateRole(role *entity.Role, operatorID uint64) error
	DeleteRole(id uint64, operatorID uint64) error
	GetRoleByID(id uint64) (*entity.Role, error)
	GetAllRoles() ([]entity.Role, error)
	BindResource(roleID, resourceID uint64, isWrite bool, operatorID uint64) error
	UnbindResource(roleID, resourceID uint64, operatorID uint64) error
	BindOrgScope(roleID, orgUnitID uint64, includeDescendants bool, operatorID uint64) error
	UnbindOrgScope(roleID, orgUnitID uint64, operatorID uint64) error
	AssignRoleToUser(userID, roleID, orgUnitID uint64, operatorID uint64) error
	RevokeRoleFromUser(userID, roleID, orgUnitID uint64, operatorID uint64) error
	GetRoleResources(roleID uint64) ([]entity.RoleResource, error)
	BindMenu(roleID, menuID uint64, operatorID uint64) error
	UnbindMenu(roleID, menuID uint64, operatorID uint64) error
	GetRoleMenus(roleID uint64) ([]entity.RoleMenu, error)
}
