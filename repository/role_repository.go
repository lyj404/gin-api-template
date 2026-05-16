package repository

import (
	"github.com/lyj404/gin-api-template/domain/entity"
	"github.com/lyj404/gin-api-template/domain/repositories"
	"github.com/lyj404/gin-api-template/global"
)

type roleRepository struct{}

func NewRoleRepository() repositories.RoleRepository {
	return &roleRepository{}
}

func (r *roleRepository) Create(role *entity.Role) error {
	return global.G_DB.Create(role).Error
}

func (r *roleRepository) Update(role *entity.Role) error {
	return global.G_DB.Save(role).Error
}

func (r *roleRepository) Delete(id uint64) error {
	return global.G_DB.Delete(&entity.Role{}, id).Error
}

func (r *roleRepository) GetByID(id uint64) (*entity.Role, error) {
	var role entity.Role
	err := global.G_DB.First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) GetAll() ([]entity.Role, error) {
	var roles []entity.Role
	err := global.G_DB.Find(&roles).Error
	return roles, err
}

func (r *roleRepository) BindResource(roleID, resourceID uint64, isWrite bool) error {
	roleResource := entity.RoleResource{
		RoleID:     roleID,
		ResourceID: resourceID,
		IsRead:     true,
		IsWrite:    isWrite,
	}
	return global.G_DB.Create(&roleResource).Error
}

func (r *roleRepository) UnbindResource(roleID, resourceID uint64) error {
	return global.G_DB.Where("role_id = ? AND resource_id = ?", roleID, resourceID).Delete(&entity.RoleResource{}).Error
}

func (r *roleRepository) GetRoleResources(roleID uint64) ([]entity.RoleResource, error) {
	var resources []entity.RoleResource
	err := global.G_DB.Preload("Resource").Where("role_id = ?", roleID).Find(&resources).Error
	return resources, err
}

func (r *roleRepository) BindOrgScope(roleID, orgUnitID uint64, includeDescendants bool) error {
	roleOrgScope := entity.RoleOrgScope{
		RoleID:             roleID,
		OrgUnitID:          orgUnitID,
		IncludeDescendants: includeDescendants,
	}
	return global.G_DB.Create(&roleOrgScope).Error
}

func (r *roleRepository) UnbindOrgScope(roleID, orgUnitID uint64) error {
	return global.G_DB.Where("role_id = ? AND org_unit_id = ?", roleID, orgUnitID).Delete(&entity.RoleOrgScope{}).Error
}

func (r *roleRepository) GetRoleOrgScopes(roleID uint64) ([]entity.RoleOrgScope, error) {
	var scopes []entity.RoleOrgScope
	err := global.G_DB.Preload("OrgUnit").Where("role_id = ?", roleID).Find(&scopes).Error
	return scopes, err
}

func (r *roleRepository) BindMenu(roleID, menuID uint64) error {
	rm := entity.RoleMenu{RoleID: roleID, MenuID: menuID}
	return global.G_DB.Create(&rm).Error
}

func (r *roleRepository) UnbindMenu(roleID, menuID uint64) error {
	return global.G_DB.Where("role_id = ? AND menu_id = ?", roleID, menuID).Delete(&entity.RoleMenu{}).Error
}

func (r *roleRepository) GetRoleMenus(roleID uint64) ([]entity.RoleMenu, error) {
	var menus []entity.RoleMenu
	err := global.G_DB.Preload("Menu").Where("role_id = ?", roleID).Find(&menus).Error
	return menus, err
}
