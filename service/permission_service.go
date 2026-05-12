package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/lyj404/gin-api-template/config"
	"github.com/lyj404/gin-api-template/domain/entity"
	"github.com/lyj404/gin-api-template/domain/services"
	"github.com/lyj404/gin-api-template/global"
)

// permissionServiceImpl 权限服务实现
type permissionServiceImpl struct{}

// NewPermissionService 创建权限服务实例
func NewPermissionService() services.PermissionService {
	return &permissionServiceImpl{}
}

// CheckPermission 检查用户是否有访问指定资源的权限
func (s *permissionServiceImpl) CheckPermission(userID uint, resource string, method string) (bool, error) {
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

// CheckEntityPermission 检查用户是否有操作指定实体的权限
func (s *permissionServiceImpl) CheckEntityPermission(userID uint, entityType string, entityID uint, action string) (bool, error) {
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

// GetUserPermissions 获取用户的权限列表
func (s *permissionServiceImpl) GetUserPermissions(userID uint) ([]services.PermissionInfo, error) {
	return s.getUserPermissions(userID)
}

// GetUserOrgScope 获取用户的组织范围
func (s *permissionServiceImpl) GetUserOrgScope(userID uint) ([]services.OrgScopeInfo, error) {
	return s.getUserOrgScope(userID)
}

// ClearUserCache 清除用户权限缓存
func (s *permissionServiceImpl) ClearUserCache(userID uint) error {
	if !config.CfgRedis.Enabled {
		return nil
	}

	cacheKey := fmt.Sprintf("user:permissions:%d", userID)
	return global.G_REDIS.Del(context.Background(), cacheKey).Err()
}

// GetUserMenus 获取用户可见的菜单树（根据用户权限过滤）
func (s *permissionServiceImpl) GetUserMenus(userID uint) ([]services.MenuTreeNode, error) {
	// 获取用户权限
	permissions, err := s.getUserPermissions(userID)
	if err != nil {
		return nil, err
	}

	// 获取用户可见的资源名列表
	resourceSet := make(map[string]bool)
	for _, perm := range permissions {
		resourceSet[perm.ResourceName] = true
	}

	// 获取所有菜单
	var menus []entity.Menu
	err = global.G_DB.Preload("Resource").Order("order_num ASC").Find(&menus).Error
	if err != nil {
		return nil, err
	}

	// 过滤菜单：只保留用户有权限且可见的菜单
	visibleMenus := make([]entity.Menu, 0)
	for _, menu := range menus {
		if !menu.IsVisible || menu.Status != "enabled" {
			continue
		}
		if menu.Resource != nil && resourceSet[menu.Resource.Name] {
			visibleMenus = append(visibleMenus, menu)
		}
	}

	// 构建菜单树
	return s.buildMenuTree(visibleMenus), nil
}

// buildMenuTree 将扁平菜单列表构建为树形结构
func (s *permissionServiceImpl) buildMenuTree(menus []entity.Menu) []services.MenuTreeNode {
	menuMap := make(map[uint]*services.MenuTreeNode)
	var roots []services.MenuTreeNode

	// 第一遍：创建所有节点映射
	for _, menu := range menus {
		menuMap[menu.ID] = &services.MenuTreeNode{
			ID:        menu.ID,
			Name:      menu.Name,
			Path:      menu.Path,
			Component: menu.Component,
			Icon:      menu.Icon,
			OrderNum:  menu.OrderNum,
		}
	}

	// 第二遍：建立父子关系
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

// getUserPermissions 获取用户权限（内部方法，支持缓存）
func (s *permissionServiceImpl) getUserPermissions(userID uint) ([]services.PermissionInfo, error) {
	cacheKey := fmt.Sprintf("user:permissions:%d", userID)

	// 如果启用Redis，尝试从缓存获取
	if config.CfgRedis.Enabled {
		cached, err := global.G_REDIS.Get(context.Background(), cacheKey).Result()
		if err == nil && cached != "" {
			var permissions []services.PermissionInfo
			if err := json.Unmarshal([]byte(cached), &permissions); err == nil {
				return permissions, nil
			}
		}
	}

	// 从数据库获取用户角色和权限
	var userRoles []entity.UserRole
	err := global.G_DB.Preload("Role.RoleResources.Resource").Where("user_id = ?", userID).Find(&userRoles).Error
	if err != nil {
		return nil, err
	}

	// 合并同一资源的权限
	permissionMap := make(map[string]*services.PermissionInfo)
	for _, userRole := range userRoles {
		for _, roleResource := range userRole.Role.RoleResources {
			if perm, exists := permissionMap[roleResource.Resource.Name]; exists {
				perm.IsRead = perm.IsRead || roleResource.IsRead
				perm.IsWrite = perm.IsWrite || roleResource.IsWrite
			} else {
				permissionMap[roleResource.Resource.Name] = &services.PermissionInfo{
					ResourceName: roleResource.Resource.Name,
					IsRead:       roleResource.IsRead,
					IsWrite:      roleResource.IsWrite,
				}
			}
		}
	}

	permissions := make([]services.PermissionInfo, 0, len(permissionMap))
	for _, perm := range permissionMap {
		permissions = append(permissions, *perm)
	}

	// 如果启用Redis，缓存权限
	if config.CfgRedis.Enabled {
		data, _ := json.Marshal(permissions)
		global.G_REDIS.Set(context.Background(), cacheKey, data, 0)
	}

	return permissions, nil
}

// getUserOrgScope 获取用户组织范围（内部方法）
func (s *permissionServiceImpl) getUserOrgScope(userID uint) ([]services.OrgScopeInfo, error) {
	var userRoles []entity.UserRole
	err := global.G_DB.Preload("Role.RoleOrgScopes.OrgUnit").Where("user_id = ?", userID).Find(&userRoles).Error
	if err != nil {
		return nil, err
	}

	scopeMap := make(map[uint]*services.OrgScopeInfo)
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

// matchPattern 通配符模式匹配
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

// isWriteMethod 判断是否为写操作方法
func (s *permissionServiceImpl) isWriteMethod(method string) bool {
	return method != "GET" && method != "HEAD"
}
