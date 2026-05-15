package repositories

import (
	"github.com/lyj404/gin-api-template/domain/entity"
)

type RoleRepository interface {
	Create(role *entity.Role) error
	Update(role *entity.Role) error
	Delete(id uint) error
	GetByID(id uint) (*entity.Role, error)
	GetAll() ([]entity.Role, error)
	BindResource(roleID, resourceID uint, isWrite bool) error
	UnbindResource(roleID, resourceID uint) error
	GetRoleResources(roleID uint) ([]entity.RoleResource, error)
	BindOrgScope(roleID, orgUnitID uint, includeDescendants bool) error
	UnbindOrgScope(roleID, orgUnitID uint) error
	GetRoleOrgScopes(roleID uint) ([]entity.RoleOrgScope, error)
	BindMenu(roleID, menuID uint) error
	UnbindMenu(roleID, menuID uint) error
	GetRoleMenus(roleID uint) ([]entity.RoleMenu, error)
}
