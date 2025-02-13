package dto

// RefreshTokenRequest 刷新token的请求
// @Description 刷新token时请求的参数
type RefreshTokenRequest struct {
	// @Description 刷新token
	// @Required
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// RefreshTokenResponse 刷新token的响应
// @Description 刷新token成功之后的响应数据
type RefreshTokenResponse struct {
	// @Description 访问token
	AccessToken string `json:"accessToken"`
	// @Description 刷新token
	RefreshToken string `json:"refreshToken"`
}
