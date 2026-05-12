package services

// MenuTreeNode 菜单树节点，用于返回给前端的用户菜单
type MenuTreeNode struct {
	ID        uint           `json:"id"`
	Name      string         `json:"name"`
	Path      string         `json:"path"`
	Component string         `json:"component"`
	Icon      string         `json:"icon"`
	OrderNum  int            `json:"order_num"`
	Children  []MenuTreeNode `json:"children,omitempty"`
}

// PermissionService 权限服务接口，定义权限检查相关的业务逻辑
type PermissionService interface {
	// CheckPermission 检查用户是否有访问指定资源的权限
	CheckPermission(userID uint, resource string, method string) (bool, error)

	// CheckEntityPermission 检查用户是否有操作指定实体的权限
	CheckEntityPermission(userID uint, entityType string, entityID uint, action string) (bool, error)

	// GetUserPermissions 获取用户的权限列表
	GetUserPermissions(userID uint) ([]PermissionInfo, error)

	// GetUserOrgScope 获取用户的组织范围
	GetUserOrgScope(userID uint) ([]OrgScopeInfo, error)

	// ClearUserCache 清除用户权限缓存
	ClearUserCache(userID uint) error

	// GetUserMenus 获取用户可见的菜单树（根据用户权限过滤）
	GetUserMenus(userID uint) ([]MenuTreeNode, error)
}

// PermissionInfo 权限信息结构
type PermissionInfo struct {
	ResourceName string `json:"resource_name"`
	IsRead       bool   `json:"is_read"`
	IsWrite      bool   `json:"is_write"`
}

// OrgScopeInfo 组织范围信息结构
type OrgScopeInfo struct {
	OrgUnitID          uint   `json:"org_unit_id"`
	IncludeDescendants bool   `json:"include_descendants"`
	Path               string `json:"path"`
}
