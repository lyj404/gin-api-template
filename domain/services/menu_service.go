package services

import "github.com/lyj404/gin-api-template/domain/entity"

type MenuService interface {
	CreateMenu(menu *entity.Menu, operatorID uint) error
	UpdateMenu(menu *entity.Menu, operatorID uint) error
	DeleteMenu(id uint, operatorID uint) error
	GetMenuByID(id uint) (*entity.Menu, error)
	GetAllMenus() ([]entity.Menu, error)
	GetMenuTree() ([]entity.Menu, error)
	BindResource(menuID, resourceID uint, operatorID uint) error
	UnbindResource(menuID, resourceID uint, operatorID uint) error
	GetMenuResources(menuID uint) ([]entity.MenuResource, error)
}
