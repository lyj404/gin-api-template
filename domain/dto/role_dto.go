package dto

type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateRoleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RoleResponse struct {
	ID          uint64   `json:"id,string"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsSystem    bool   `json:"is_system"`
}

type RoleDetailResponse struct {
	ID          uint64                   `json:"id,string"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	IsSystem    bool                   `json:"is_system"`
	Resources   []RoleResourceResponse `json:"resources,omitempty"`
	Menus       []RoleMenuResponse     `json:"menus,omitempty"`
}

type RoleMenuResponse struct {
	ID         uint64            `json:"id,string"`
	RoleID     uint64            `json:"role_id,string"`
	MenuID     uint64            `json:"menu_id,string"`
	Menu       *MenuBriefResponse `json:"menu,omitempty"`
}

type MenuBriefResponse struct {
	ID     uint64   `json:"id,string"`
	Name   string `json:"name"`
	Path   string `json:"path"`
	Icon   string `json:"icon"`
}

type BindRoleResourceRequest struct {
	ResourceID uint64 `json:"resource_id,string" binding:"required"`
	IsWrite    bool `json:"is_write"`
}

type BindRoleMenuRequest struct {
	MenuID uint64 `json:"menu_id,string" binding:"required"`
}

type RoleResourceResponse struct {
	ID         uint64                  `json:"id,string"`
	RoleID     uint64                  `json:"role_id,string"`
	ResourceID uint64                  `json:"resource_id,string"`
	IsRead     bool                  `json:"is_read"`
	IsWrite    bool                  `json:"is_write"`
	Resource   *ResourceBriefResponse `json:"resource,omitempty"`
}

type ResourceBriefResponse struct {
	ID          uint64   `json:"id,string"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Pattern     string `json:"pattern"`
	Method      string `json:"method"`
	Entity      string `json:"entity"`
	Action      string `json:"action"`
	Description string `json:"description"`
}
