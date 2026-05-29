package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/lyj404/gin-api-template/config"
	"github.com/lyj404/gin-api-template/domain/entity"
	"github.com/lyj404/gin-api-template/domain/services"
	"github.com/lyj404/gin-api-template/global"
)

type permissionServiceImpl struct{}

func NewPermissionService() services.PermissionService {
	return &permissionServiceImpl{}
}

func (s *permissionServiceImpl) CheckPermission(userID uint64, resource string, method string) (bool, error) {
	permissions, err := s.getUserPermissions(userID)
	if err != nil {
		return false, err
	}

	for _, perm := range permissions {
		if s.matchPattern(perm.ResourceName, resource) {
			isWrite := s.isWriteMethod(method)
			if isWrite && perm.IsWrite {
				return true, nil
			} else if !isWrite && perm.IsRead {
				return true, nil
			}
		}
	}

	return false, nil
}

func (s *permissionServiceImpl) CheckEntityPermission(userID uint64, entityType string, entityID uint64, action string) (bool, error) {
	entityResourceName := fmt.Sprintf("entity:%s:%s", entityType, action)

	permissions, err := s.getUserPermissions(userID)
	if err != nil {
		return false, err
	}

	hasPermission := false
	for _, perm := range permissions {
		if perm.ResourceName == "entity:all" || s.matchPattern(perm.ResourceName, entityResourceName) {
			hasPermission = true
			break
		}
	}

	if !hasPermission {
		return false, nil
	}

	orgScope, err := s.getUserOrgScope(userID)
	if err != nil {
		return false, err
	}

	if len(orgScope) == 0 {
		return false, nil
	}

	var binding entity.OrgEntityBinding
	err = global.G_DB.Where("entity_type = ? AND entity_id = ?", entityType, entityID).First(&binding).Error
	if err != nil {
		return false, nil
	}

	for _, scope := range orgScope {
		if scope.OrgUnitID == binding.OrgUnitID {
			return true, nil
		}

		if scope.IncludeDescendants {
			var orgUnit entity.OrgUnit
			if err := global.G_DB.First(&orgUnit, binding.OrgUnitID).Error; err == nil {
				if strings.HasPrefix(orgUnit.Path, scope.Path) {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

func (s *permissionServiceImpl) GetUserPermissions(userID uint64) ([]services.PermissionInfo, error) {
	return s.getUserPermissions(userID)
}

func (s *permissionServiceImpl) GetUserOrgScope(userID uint64) ([]services.OrgScopeInfo, error) {
	return s.getUserOrgScope(userID)
}

func (s *permissionServiceImpl) ClearUserCache(userID uint64) error {
	if !config.CfgRedis.Enabled {
		return nil
	}

	cacheKey := fmt.Sprintf("user:permissions:%d", userID)
	return global.G_REDIS.Del(context.Background(), cacheKey).Err()
}

func (s *permissionServiceImpl) GetUserMenus(userID uint64) ([]services.MenuTreeNode, error) {
	var userRoles []entity.UserRole
	err := global.G_DB.Preload("Role.RoleMenus.Menu.Resources").Where("user_id = ?", userID).Find(&userRoles).Error
	if err != nil {
		return nil, err
	}

	menuMap := make(map[uint64]*entity.Menu)
	for _, ur := range userRoles {
		for _, rm := range ur.Role.RoleMenus {
			if rm.Menu != nil && rm.Menu.IsVisible && rm.Menu.Status == "enabled" {
				if _, exists := menuMap[rm.Menu.ID]; !exists {
					menuMap[rm.Menu.ID] = rm.Menu
				}
			}
		}
	}

	allMenus := make([]entity.Menu, 0, len(menuMap))
	for _, m := range menuMap {
		allMenus = append(allMenus, *m)
	}
	sort.Slice(allMenus, func(i, j int) bool {
		return allMenus[i].OrderNum < allMenus[j].OrderNum
	})

	return s.buildMenuTree(allMenus), nil
}

func (s *permissionServiceImpl) buildMenuTree(menus []entity.Menu) []services.MenuTreeNode {
	menuMap := make(map[uint64]*services.MenuTreeNode)
	var roots []services.MenuTreeNode

	for _, menu := range menus {
		menuMap[menu.ID] = &services.MenuTreeNode{
			ID:       menu.ID,
			Name:     menu.Name,
			Path:     menu.Path,
			Icon:     menu.Icon,
			OrderNum: menu.OrderNum,
		}
	}

	for _, menu := range menus {
		if menu.ParentID != nil {
			if parent, exists := menuMap[*menu.ParentID]; exists {
				parent.Children = append(parent.Children, *menuMap[menu.ID])
			}
		} else {
			roots = append(roots, *menuMap[menu.ID])
		}
	}

	return roots
}

func (s *permissionServiceImpl) getUserPermissions(userID uint64) ([]services.PermissionInfo, error) {
	cacheKey := fmt.Sprintf("user:permissions:%d", userID)

	if config.CfgRedis.Enabled {
		cached, err := global.G_REDIS.Get(context.Background(), cacheKey).Result()
		if err == nil && cached != "" {
			var permissions []services.PermissionInfo
			if err := json.Unmarshal([]byte(cached), &permissions); err == nil {
				return permissions, nil
			}
		}
	}

	var userRoles []entity.UserRole
	err := global.G_DB.
		Preload("Role.RoleResources.Resource").
		Preload("Role.RoleMenus.Menu.Resources").
		Where("user_id = ?", userID).
		Find(&userRoles).Error
	if err != nil {
		return nil, err
	}

	permissionMap := make(map[string]*services.PermissionInfo)

	// 1. 收集角色直接绑定的资源
	for _, userRole := range userRoles {
		for _, roleResource := range userRole.Role.RoleResources {
			s.mergePermission(permissionMap, roleResource.Resource.Name, roleResource.IsRead, roleResource.IsWrite)
		}
	}

	// 2. 收集角色通过菜单绑定的资源
	for _, userRole := range userRoles {
		for _, roleMenu := range userRole.Role.RoleMenus {
			if roleMenu.Menu != nil {
				for _, res := range roleMenu.Menu.Resources {
					s.mergePermission(permissionMap, res.Name, true, false)
				}
			}
		}
	}

	permissions := make([]services.PermissionInfo, 0, len(permissionMap))
	for _, perm := range permissionMap {
		permissions = append(permissions, *perm)
	}

	if config.CfgRedis.Enabled {
		data, _ := json.Marshal(permissions)
		global.G_REDIS.Set(context.Background(), cacheKey, data, 0)
	}

	return permissions, nil
}

func (s *permissionServiceImpl) mergePermission(permMap map[string]*services.PermissionInfo, name string, isRead, isWrite bool) {
	if perm, exists := permMap[name]; exists {
		perm.IsRead = perm.IsRead || isRead
		perm.IsWrite = perm.IsWrite || isWrite
	} else {
		permMap[name] = &services.PermissionInfo{
			ResourceName: name,
			IsRead:       isRead,
			IsWrite:      isWrite,
		}
	}
}

func (s *permissionServiceImpl) getUserOrgScope(userID uint64) ([]services.OrgScopeInfo, error) {
	var userRoles []entity.UserRole
	err := global.G_DB.Preload("Role.RoleOrgScopes.OrgUnit").Where("user_id = ?", userID).Find(&userRoles).Error
	if err != nil {
		return nil, err
	}

	scopeMap := make(map[uint64]*services.OrgScopeInfo)
	for _, userRole := range userRoles {
		for _, scope := range userRole.Role.RoleOrgScopes {
			if info, exists := scopeMap[scope.OrgUnitID]; exists {
				info.IncludeDescendants = info.IncludeDescendants || scope.IncludeDescendants
			} else {
				scopeMap[scope.OrgUnitID] = &services.OrgScopeInfo{
					OrgUnitID:          scope.OrgUnit.ID,
					IncludeDescendants: scope.IncludeDescendants,
					Path:               scope.OrgUnit.Path,
				}
			}
		}
	}

	scopes := make([]services.OrgScopeInfo, 0, len(scopeMap))
	for _, scope := range scopeMap {
		scopes = append(scopes, *scope)
	}

	return scopes, nil
}

func (s *permissionServiceImpl) HasSystemRole(userID uint64) (bool, error) {
	var count int64
	err := global.G_DB.Model(&entity.UserRole{}).
		Joins(`JOIN role ON role.id = user_role.role_id AND role.deleted_at IS NULL`).
		Where("user_role.user_id = ? AND role.is_system = ?", userID, true).
		Count(&count).Error
	return count > 0, err
}

func (s *permissionServiceImpl) matchPattern(pattern, target string) bool {
	if pattern == "*" {
		return true
	}

	if strings.HasSuffix(pattern, "*") {
		prefix := strings.TrimSuffix(pattern, "*")
		return strings.HasPrefix(target, prefix)
	}

	return pattern == target
}

func (s *permissionServiceImpl) isWriteMethod(method string) bool {
	return method != "GET" && method != "HEAD"
}
