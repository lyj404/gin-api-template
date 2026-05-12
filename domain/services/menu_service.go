package services

import "github.com/lyj404/gin-api-template/domain/entity"

// MenuService 菜单服务接口，定义菜单管理的业务逻辑
type MenuService interface {
	// CreateMenu 创建菜单
	CreateMenu(menu *entity.Menu, operatorID uint) error

	// UpdateMenu 更新菜单
	UpdateMenu(menu *entity.Menu, operatorID uint) error

	// DeleteMenu 删除菜单
	DeleteMenu(id uint, operatorID uint) error

	// GetMenuByID 根据ID获取菜单
	GetMenuByID(id uint) (*entity.Menu, error)

	// GetAllMenus 获取所有菜单（扁平结构）
	GetAllMenus() ([]entity.Menu, error)

	// GetMenuTree 获取菜单树形结构
	GetMenuTree() ([]entity.Menu, error)
}