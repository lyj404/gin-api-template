package services

import "github.com/lyj404/gin-api-template/domain/entity"

type MenuService interface {
	CreateMenu(menu *entity.Menu, operatorID uint64) error
	UpdateMenu(menu *entity.Menu, operatorID uint64) error
	DeleteMenu(id uint64, operatorID uint64) error
	GetMenuByID(id uint64) (*entity.Menu, error)
	GetAllMenus() ([]entity.Menu, error)
	GetMenuTree() ([]entity.Menu, error)
	BindResource(menuID, resourceID uint64, operatorID uint64) error
	UnbindResource(menuID, resourceID uint64, operatorID uint64) error
	GetMenuResources(menuID uint64) ([]entity.MenuResource, error)
}
