package dto

// LoginRequest 登录请求
// @Description 登录时请求的参数
type LoginRequest struct {
	// @Description 用户邮箱
	// @Required
	Email string `json:"email" binding:"required,email"`
	// @Description 用户密码
	// @Required
	Password string `json:"password" binding:"required,min=6"`
}

// LoginResponse 登录响应
// @Description 登录成功后的响应数据
type LoginResponse struct {
	// @Description 访问令牌
	AccessToken string `json:"accessToken"`
	// @Description 刷新令牌
	RefreshToken string `json:"refreshToken"`
}
