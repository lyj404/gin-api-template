package repository

import (
	"github.com/lyj404/gin-api-template/domain"
	"github.com/lyj404/gin-api-template/domain/entity"
	"github.com/lyj404/gin-api-template/global"
)

// menuRepository 菜单仓储实现
type menuRepository struct{}

// NewMenuRepository 创建菜单仓储实例
func NewMenuRepository() domain.MenuRepository {
	return &menuRepository{}
}

// Create 创建菜单
func (r *menuRepository) Create(menu *entity.Menu) error {
	return global.G_DB.Create(menu).Error
}

// Update 更新菜单
func (r *menuRepository) Update(menu *entity.Menu) error {
	return global.G_DB.Save(menu).Error
}

// Delete 删除菜单（软删除）
func (r *menuRepository) Delete(id uint) error {
	return global.G_DB.Delete(&entity.Menu{}, id).Error
}

// GetByID 根据ID获取菜单
func (r *menuRepository) GetByID(id uint) (*entity.Menu, error) {
	var menu entity.Menu
	err := global.G_DB.Preload("Resource").First(&menu, id).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

// GetAll 获取所有菜单（包含资源信息）
func (r *menuRepository) GetAll() ([]entity.Menu, error) {
	var menus []entity.Menu
	err := global.G_DB.Preload("Resource").Order("order_num ASC").Find(&menus).Error
	return menus, err
}

// GetByParentID 根据父ID获取子菜单
func (r *menuRepository) GetByParentID(parentID *uint) ([]entity.Menu, error) {
	var menus []entity.Menu
	query := global.G_DB.Preload("Resource").Order("order_num ASC")
	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", *parentID)
	}
	err := query.Find(&menus).Error
	return menus, err
}

// GetRootMenus 获取根菜单
func (r *menuRepository) GetRootMenus() ([]entity.Menu, error) {
	var menus []entity.Menu
	err := global.G_DB.Preload("Resource").Where("parent_id IS NULL").Order("order_num ASC").Find(&menus).Error
	return menus, err
}

// GetByResourceID 根据资源ID获取菜单
func (r *menuRepository) GetByResourceID(resourceID uint) (*entity.Menu, error) {
	var menu entity.Menu
	err := global.G_DB.Where("resource_id = ?", resourceID).First(&menu).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

// HasChildren 检查菜单是否有子菜单
func (r *menuRepository) HasChildren(id uint) (bool, error) {
	var count int64
	err := global.G_DB.Model(&entity.Menu{}).Where("parent_id = ?", id).Count(&count).Error
	return count > 0, err
}