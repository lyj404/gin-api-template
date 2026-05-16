package repositories

import (
	"github.com/lyj404/gin-api-template/domain/entity"
)

type RoleRepository interface {
	Create(role *entity.Role) error
	Update(role *entity.Role) error
	Delete(id uint64) error
	GetByID(id uint64) (*entity.Role, error)
	GetAll() ([]entity.Role, error)
	BindResource(roleID, resourceID uint64, isWrite bool) error
	UnbindResource(roleID, resourceID uint64) error
	GetRoleResources(roleID uint64) ([]entity.RoleResource, error)
	BindOrgScope(roleID, orgUnitID uint64, includeDescendants bool) error
	UnbindOrgScope(roleID, orgUnitID uint64) error
	GetRoleOrgScopes(roleID uint64) ([]entity.RoleOrgScope, error)
	BindMenu(roleID, menuID uint64) error
	UnbindMenu(roleID, menuID uint64) error
	GetRoleMenus(roleID uint64) ([]entity.RoleMenu, error)
}
