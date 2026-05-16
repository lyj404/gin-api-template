package service

import (
	"encoding/json"
	"fmt"

	"github.com/lyj404/gin-api-template/domain"
	"github.com/lyj404/gin-api-template/domain/entity"
	"github.com/lyj404/gin-api-template/domain/services"
	"github.com/lyj404/gin-api-template/global"
	"gorm.io/gorm"
)

// menuServiceImpl 菜单服务实现
type menuServiceImpl struct {
	menuRepo domain.MenuRepository
}

// NewMenuService 创建菜单服务实例
func NewMenuService(menuRepo domain.MenuRepository) services.MenuService {
	return &menuServiceImpl{
		menuRepo: menuRepo,
	}
}

// CreateMenu 创建菜单
func (s *menuServiceImpl) CreateMenu(menu *entity.Menu, operatorID uint64) error {
	return global.G_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(menu).Error; err != nil {
			return err
		}

		menuJSON, _ := json.Marshal(menu)
		description := fmt.Sprintf("创建菜单: %s", menu.Name)
		return s.createAuditLog(tx, operatorID, "create", "menu", menu.ID, "", string(menuJSON), description)
	})
}

// UpdateMenu 更新菜单
func (s *menuServiceImpl) UpdateMenu(menu *entity.Menu, operatorID uint64) error {
	return global.G_DB.Transaction(func(tx *gorm.DB) error {
		oldMenu, err := s.menuRepo.GetByID(menu.ID)
		if err != nil {
			return err
		}

		if err := tx.Save(menu).Error; err != nil {
			return err
		}

		oldJSON, _ := json.Marshal(oldMenu)
		newJSON, _ := json.Marshal(menu)
		description := fmt.Sprintf("更新菜单: %s", menu.Name)
		return s.createAuditLog(tx, operatorID, "update", "menu", menu.ID, string(oldJSON), string(newJSON), description)
	})
}

// DeleteMenu 删除菜单
func (s *menuServiceImpl) DeleteMenu(id uint64, operatorID uint64) error {
	return global.G_DB.Transaction(func(tx *gorm.DB) error {
		menu, err := s.menuRepo.GetByID(id)
		if err != nil {
			return err
		}

		if err := tx.Delete(&entity.Menu{}, id).Error; err != nil {
			return err
		}

		menuJSON, _ := json.Marshal(menu)
		description := fmt.Sprintf("删除菜单: %s", menu.Name)
		return s.createAuditLog(tx, operatorID, "delete", "menu", id, string(menuJSON), "", description)
	})
}

// GetMenuByID 根据ID获取菜单
func (s *menuServiceImpl) GetMenuByID(id uint64) (*entity.Menu, error) {
	return s.menuRepo.GetByID(id)
}

// GetAllMenus 获取所有菜单（扁平结构）
func (s *menuServiceImpl) GetAllMenus() ([]entity.Menu, error) {
	return s.menuRepo.GetAll()
}

// GetMenuTree 获取菜单树形结构
func (s *menuServiceImpl) GetMenuTree() ([]entity.Menu, error) {
	allMenus, err := s.menuRepo.GetAll()
	if err != nil {
		return nil, err
	}

	return buildMenuTree(allMenus), nil
}

// buildMenuTree 将扁平菜单列表构建为树形结构
func buildMenuTree(menus []entity.Menu) []entity.Menu {
	menuMap := make(map[uint64]*entity.Menu)
	var roots []entity.Menu

	// 第一遍：创建所有节点映射
	for i := range menus {
		menus[i].Children = []entity.Menu{}
		menuMap[menus[i].ID] = &menus[i]
	}

	// 第二遍：建立父子关系
	for i := range menus {
		if menus[i].ParentID != nil {
			if parent, exists := menuMap[*menus[i].ParentID]; exists {
				parent.Children = append(parent.Children, menus[i])
			}
		} else {
			roots = append(roots, menus[i])
		}
	}

	return roots
}

// createAuditLog 创建审计日志
func (s *menuServiceImpl) BindResource(menuID, resourceID uint64, operatorID uint64) error {
	return global.G_DB.Transaction(func(tx *gorm.DB) error {
		mr := entity.MenuResource{MenuID: menuID, ResourceID: resourceID}
		if err := tx.Create(&mr).Error; err != nil {
			return err
		}
		description := fmt.Sprintf("菜单 %d 绑定资源 %d", menuID, resourceID)
		return s.createAuditLog(tx, operatorID, "bind", "menu_resource", mr.ID, "", "", description)
	})
}

func (s *menuServiceImpl) UnbindResource(menuID, resourceID uint64, operatorID uint64) error {
	return global.G_DB.Transaction(func(tx *gorm.DB) error {
		description := fmt.Sprintf("菜单 %d 解绑资源 %d", menuID, resourceID)
		if err := tx.Where("menu_id = ? AND resource_id = ?", menuID, resourceID).Delete(&entity.MenuResource{}).Error; err != nil {
			return err
		}
		return s.createAuditLog(tx, operatorID, "unbind", "menu_resource", 0, "", "", description)
	})
}

func (s *menuServiceImpl) GetMenuResources(menuID uint64) ([]entity.MenuResource, error) {
	return s.menuRepo.GetMenuResources(menuID)
}

func (s *menuServiceImpl) createAuditLog(tx *gorm.DB, operatorID uint64, action, targetType string, targetID uint64, beforeData, afterData, description string) error {
	auditLog := entity.AuditLog{
		OperatorID:   operatorID,
		OperatorName: getOperatorName(tx, operatorID),
		Action:       action,
		TargetType:   targetType,
		TargetID:     targetID,
		BeforeData:   beforeData,
		AfterData:    afterData,
		Description:  description,
	}
	return tx.Create(&auditLog).Error
}