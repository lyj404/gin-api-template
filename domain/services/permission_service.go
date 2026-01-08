package services

type PermissionService interface {
	CheckPermission(userID uint, resource string, method string) (bool, error)
	CheckEntityPermission(userID uint, entityType string, entityID uint, action string) (bool, error)
	GetUserPermissions(userID uint) ([]PermissionInfo, error)
	GetUserOrgScope(userID uint) ([]OrgScopeInfo, error)
	ClearUserCache(userID uint) error
}

type PermissionInfo struct {
	ResourceName string `json:"resource_name"`
	IsRead       bool   `json:"is_read"`
	IsWrite      bool   `json:"is_write"`
}

type OrgScopeInfo struct {
	OrgUnitID          uint   `json:"org_unit_id"`
	IncludeDescendants bool   `json:"include_descendants"`
	Path               string `json:"path"`
}
