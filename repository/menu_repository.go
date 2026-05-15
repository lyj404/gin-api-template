package repository

import (
	"github.com/lyj404/gin-api-template/domain"
	"github.com/lyj404/gin-api-template/domain/entity"
	"github.com/lyj404/gin-api-template/global"
)

type menuRepository struct{}

func NewMenuRepository() domain.MenuRepository {
	return &menuRepository{}
}

func (r *menuRepository) Create(menu *entity.Menu) error {
	return global.G_DB.Create(menu).Error
}

func (r *menuRepository) Update(menu *entity.Menu) error {
	return global.G_DB.Save(menu).Error
}

func (r *menuRepository) Delete(id uint) error {
	return global.G_DB.Delete(&entity.Menu{}, id).Error
}

func (r *menuRepository) GetByID(id uint) (*entity.Menu, error) {
	var menu entity.Menu
	err := global.G_DB.Preload("Resources").First(&menu, id).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

func (r *menuRepository) GetAll() ([]entity.Menu, error) {
	var menus []entity.Menu
	err := global.G_DB.Preload("Resources").Order("order_num ASC").Find(&menus).Error
	return menus, err
}

func (r *menuRepository) GetByParentID(parentID *uint) ([]entity.Menu, error) {
	var menus []entity.Menu
	query := global.G_DB.Preload("Resources").Order("order_num ASC")
	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", *parentID)
	}
	err := query.Find(&menus).Error
	return menus, err
}

func (r *menuRepository) GetRootMenus() ([]entity.Menu, error) {
	var menus []entity.Menu
	err := global.G_DB.Preload("Resources").Where("parent_id IS NULL").Order("order_num ASC").Find(&menus).Error
	return menus, err
}

func (r *menuRepository) HasChildren(id uint) (bool, error) {
	var count int64
	err := global.G_DB.Model(&entity.Menu{}).Where("parent_id = ?", id).Count(&count).Error
	return count > 0, err
}

func (r *menuRepository) BindResource(menuID, resourceID uint) error {
	mr := entity.MenuResource{MenuID: menuID, ResourceID: resourceID}
	return global.G_DB.Create(&mr).Error
}

func (r *menuRepository) UnbindResource(menuID, resourceID uint) error {
	return global.G_DB.Where("menu_id = ? AND resource_id = ?", menuID, resourceID).Delete(&entity.MenuResource{}).Error
}

func (r *menuRepository) GetMenuResources(menuID uint) ([]entity.MenuResource, error) {
	var res []entity.MenuResource
	err := global.G_DB.Preload("Resource").Where("menu_id = ?", menuID).Find(&res).Error
	return res, err
}
