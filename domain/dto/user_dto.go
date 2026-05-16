package dto

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Name     string   `json:"name" binding:"required"`
	Email    string   `json:"email" binding:"required,email"`
	Password string   `json:"password" binding:"required,min=6"`
	RoleIDs  []string `json:"role_ids"`
	OrgUnitID uint64  `json:"org_unit_id,string"`
}

type UpdateUserRequest struct {
	Name     string   `json:"name"`
	Email    string   `json:"email" binding:"omitempty,email"`
	Password string   `json:"password" binding:"omitempty,min=6"`
	RoleIDs  []string `json:"role_ids"`
	OrgUnitID uint64  `json:"org_unit_id,string"`
}

type UserResponse struct {
	ID      uint64   `json:"id,string"`
	Name    string   `json:"name"`
	Email   string   `json:"email"`
	RoleIDs []string `json:"role_ids"`
	Roles   []string `json:"roles"`
}

type ProfileResponse struct {
	ID        uint64   `json:"id,string"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// UpdateProfileRequest 更新个人信息请求
type UpdateProfileRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}
