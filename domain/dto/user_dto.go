package dto

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	RoleIDs  []uint `json:"role_ids"`
	OrgUnitID uint  `json:"org_unit_id"`
}

// UpdateUserRequest 更新用户请求（密码与角色为可选）
type UpdateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"omitempty,min=6"`
	RoleIDs  []uint `json:"role_ids"`
	OrgUnitID uint  `json:"org_unit_id"`
}

// UserResponse 用户响应
type UserResponse struct {
	ID      uint     `json:"id"`
	Name    string   `json:"name"`
	Email   string   `json:"email"`
	RoleIDs []uint   `json:"role_ids"`
	Roles   []string `json:"roles"`
}
