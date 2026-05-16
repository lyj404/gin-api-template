package dto

type CreateMenuRequest struct {
	Name       string  `json:"name" binding:"required"`
	ParentID   *uint64   `json:"parent_id,string"`
	Path       string  `json:"path"`
	Icon       string  `json:"icon"`
	OrderNum   int     `json:"order_num"`
	IsVisible  bool    `json:"is_visible"`
}

type UpdateMenuRequest struct {
	Name       string  `json:"name"`
	ParentID   *uint64   `json:"parent_id,string"`
	Path       string  `json:"path"`
	Icon       string  `json:"icon"`
	OrderNum   int     `json:"order_num"`
	IsVisible  *bool   `json:"is_visible"`
	Status     string  `json:"status"`
}

type MenuResponse struct {
	ID         uint64            `json:"id,string"`
	Name       string          `json:"name"`
	ParentID   *uint64           `json:"parent_id,string"`
	Path       string          `json:"path"`
	Icon       string          `json:"icon"`
	OrderNum   int             `json:"order_num"`
	IsVisible  bool            `json:"is_visible"`
	Status     string          `json:"status"`
	Resources  []ResourceBriefResponse `json:"resources,omitempty"`
	Children   []MenuResponse  `json:"children,omitempty"`
}

type MenuTreeNode struct {
	ID         uint64           `json:"id,string"`
	Name       string         `json:"name"`
	Path       string         `json:"path"`
	Icon       string         `json:"icon"`
	OrderNum   int            `json:"order_num"`
	IsVisible  bool           `json:"is_visible"`
	Children   []MenuTreeNode `json:"children,omitempty"`
}

type UserMenuResponse struct {
	Menus []MenuTreeNode `json:"menus"`
}

type MenuListResponse struct {
	ID           uint64                   `json:"id,string"`
	Name         string                 `json:"name"`
	ParentID     *uint64                  `json:"parent_id,string"`
	Path         string                 `json:"path"`
	Icon         string                 `json:"icon"`
	OrderNum     int                    `json:"order_num"`
	IsVisible    bool                   `json:"is_visible"`
	Status       string                 `json:"status"`
	Resources    []ResourceBriefResponse `json:"resources,omitempty"`
}

type MenuResourceResponse struct {
	ID         uint64                  `json:"id,string"`
	MenuID     uint64                  `json:"menu_id,string"`
	ResourceID uint64                  `json:"resource_id,string"`
	Resource   *ResourceBriefResponse `json:"resource,omitempty"`
}

type BindMenuResourceRequest struct {
	ResourceID uint64 `json:"resource_id,string" binding:"required"`
}
