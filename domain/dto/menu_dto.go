package dto

type CreateMenuRequest struct {
	Name       string `json:"name" binding:"required"`                         // 菜单名称（必填）
	ParentID   *uint  `json:"parent_id"`                                        // 父菜单ID
	Path       string `json:"path"`                                             // 前端路由路径
	Component  string `json:"component"`                                        // 前端组件路径
	Icon       string `json:"icon"`                                            // 菜单图标
	OrderNum   int    `json:"order_num"`                                       // 排序序号
	ResourceID uint   `json:"resource_id" binding:"required"`                   // 关联资源ID（必填）
	IsVisible  bool   `json:"is_visible"`                                      // 是否显示
}

type UpdateMenuRequest struct {
	Name       string `json:"name"`                                             // 菜单名称
	ParentID   *uint  `json:"parent_id"`                                        // 父菜单ID
	Path       string `json:"path"`                                             // 前端路由路径
	Component  string `json:"component"`                                        // 前端组件路径
	Icon       string `json:"icon"`                                             // 菜单图标
	OrderNum   int    `json:"order_num"`                                        // 排序序号
	ResourceID uint   `json:"resource_id"`                                     // 关联资源ID
	IsVisible  *bool  `json:"is_visible"`                                       // 是否显示
	Status     string `json:"status"`                                           // 状态：enabled/disabled
}

type MenuResponse struct {
	ID         uint            `json:"id"`
	Name       string          `json:"name"`
	ParentID   *uint           `json:"parent_id"`
	Path       string          `json:"path"`
	Component  string          `json:"component"`
	Icon       string          `json:"icon"`
	OrderNum   int             `json:"order_num"`
	ResourceID uint            `json:"resource_id"`
	IsVisible  bool            `json:"is_visible"`
	Status     string          `json:"status"`
	Children   []MenuResponse `json:"children,omitempty"`
}

// MenuTreeNode 菜单树节点结构，用于返回给前端
type MenuTreeNode struct {
	ID         uint            `json:"id"`
	Name       string          `json:"name"`
	Path       string          `json:"path"`
	Component  string          `json:"component"`
	Icon       string          `json:"icon"`
	OrderNum   int             `json:"order_num"`
	IsVisible  bool            `json:"is_visible"`
	Children   []MenuTreeNode  `json:"children,omitempty"`
}

// UserMenuResponse 用户菜单响应，包含用户可见的菜单树
type UserMenuResponse struct {
	Menus []MenuTreeNode `json:"menus"`
}

// MenuListResponse 菜单列表响应（平面结构）
type MenuListResponse struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	ParentID   *uint  `json:"parent_id"`
	Path       string `json:"path"`
	Component  string `json:"component"`
	Icon       string `json:"icon"`
	OrderNum   int    `json:"order_num"`
	ResourceID uint   `json:"resource_id"`
	ResourceName string `json:"resource_name"` // 资源名称
	IsVisible  bool   `json:"is_visible"`
	Status     string `json:"status"`
}