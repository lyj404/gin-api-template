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
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsSystem    bool   `json:"is_system"`
}

type RoleDetailResponse struct {
	ID          uint                   `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	IsSystem    bool                   `json:"is_system"`
	Resources   []RoleResourceResponse `json:"resources,omitempty"`
}

type BindRoleResourceRequest struct {
	ResourceID uint `json:"resource_id" binding:"required"`
	IsWrite    bool `json:"is_write"`
}

type RoleResourceResponse struct {
	ID         uint                  `json:"id"`
	RoleID     uint                  `json:"role_id"`
	ResourceID uint                  `json:"resource_id"`
	IsRead     bool                  `json:"is_read"`
	IsWrite    bool                  `json:"is_write"`
	Resource   *ResourceBriefResponse `json:"resource,omitempty"`
}

type ResourceBriefResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Pattern     string `json:"pattern"`
	Method      string `json:"method"`
	Entity      string `json:"entity"`
	Action      string `json:"action"`
	Description string `json:"description"`
}
