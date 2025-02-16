package dto

// SignupRequest 注册请求
// @Description 注册时请求的参数
type SignupRequest struct {
	// @Description 用户名
	// @Required
	Name string `json:"name" binding:"required,min=3"`
	// @Description 邮箱
	// @Required
	Email string `json:"email" binding:"required,email"`
	// @Description 密码
	// @Required
	Password string `json:"password" binding:"required,min=6"`
}

// SignupResponse 注册成功响应
// @Description 注册成功后的响应数据
type SignupResponse struct {
	// @Description 访问token
	AccessToken string `json:"accessToken"`
	// @Description 刷新token
	RefreshToken string `json:"refreshToken"`
}
