package domain

import (
	"context"

	"github.com/lyj404/gin-api-template/domain/entity"
)

// UserRepo 用户仓储接口
type UserRepo interface {
	Create(c context.Context, user *entity.User) error
	GetByEmail(c context.Context, email string) (entity.User, error)
	GetByID(c context.Context, id string) (entity.User, error)
}

// MenuRepository 菜单仓储接口
type MenuRepository interface {
	Create(menu *entity.Menu) error
	Update(menu *entity.Menu) error
	Delete(id uint64) error
	GetByID(id uint64) (*entity.Menu, error)
	GetAll() ([]entity.Menu, error)
	GetByParentID(parentID *uint64) ([]entity.Menu, error)
	GetRootMenus() ([]entity.Menu, error)
	HasChildren(id uint64) (bool, error)
	BindResource(menuID, resourceID uint64) error
	UnbindResource(menuID, resourceID uint64) error
	GetMenuResources(menuID uint64) ([]entity.MenuResource, error)
}
