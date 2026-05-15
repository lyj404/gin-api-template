package dto

type CreateMenuRequest struct {
	Name       string  `json:"name" binding:"required"`
	ParentID   *uint   `json:"parent_id"`
	Path       string  `json:"path"`
	Icon       string  `json:"icon"`
	OrderNum   int     `json:"order_num"`
	IsVisible  bool    `json:"is_visible"`
}

type UpdateMenuRequest struct {
	Name       string  `json:"name"`
	ParentID   *uint   `json:"parent_id"`
	Path       string  `json:"path"`
	Icon       string  `json:"icon"`
	OrderNum   int     `json:"order_num"`
	IsVisible  *bool   `json:"is_visible"`
	Status     string  `json:"status"`
}

type MenuResponse struct {
	ID         uint            `json:"id"`
	Name       string          `json:"name"`
	ParentID   *uint           `json:"parent_id"`
	Path       string          `json:"path"`
	Icon       string          `json:"icon"`
	OrderNum   int             `json:"order_num"`
	IsVisible  bool            `json:"is_visible"`
	Status     string          `json:"status"`
	Resources  []ResourceBriefResponse `json:"resources,omitempty"`
	Children   []MenuResponse  `json:"children,omitempty"`
}

type MenuTreeNode struct {
	ID         uint           `json:"id"`
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
	ID           uint                   `json:"id"`
	Name         string                 `json:"name"`
	ParentID     *uint                  `json:"parent_id"`
	Path         string                 `json:"path"`
	Icon         string                 `json:"icon"`
	OrderNum     int                    `json:"order_num"`
	IsVisible    bool                   `json:"is_visible"`
	Status       string                 `json:"status"`
	Resources    []ResourceBriefResponse `json:"resources,omitempty"`
}

type MenuResourceResponse struct {
	ID         uint                  `json:"id"`
	MenuID     uint                  `json:"menu_id"`
	ResourceID uint                  `json:"resource_id"`
	Resource   *ResourceBriefResponse `json:"resource,omitempty"`
}

type BindMenuResourceRequest struct {
	ResourceID uint `json:"resource_id" binding:"required"`
}
