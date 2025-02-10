package handler

import (
	"gin-api-template/config"
	"gin-api-template/domain"
	"gin-api-template/domain/dto"
	"gin-api-template/domain/result"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RefreshTokenHandler struct {
	RefreshTokenService domain.RefreshTokenService
}

func (rtc *RefreshTokenHandler) RefreshToken(c *gin.Context) {
	var request dto.RefreshTokenRequest

	// 获取请求参数
	err := c.ShouldBind(&request)
	if err != nil {
		result.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// 通过token获取用户id
	id, err := rtc.RefreshTokenService.ExtractIDFromToken(request.RefreshToken, config.CfgToken.RefreshTokenSecret)
	if err != nil {
		result.ErrorResponse(c, http.StatusUnauthorized, "User not found")
		return
	}

	// 通过ID获取用户信息
	user, err := rtc.RefreshTokenService.GetUserByID(c, id)
	if err != nil {
		result.ErrorResponse(c, http.StatusUnauthorized, "User not found")
		return
	}

	// 创建新的访问token
	accessToken, err := rtc.RefreshTokenService.CreateAccessToken(&user, config.CfgToken.AccessTokenSecret, config.CfgToken.AccessTokenExpiryHour)
	if err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 创建新的刷新token
	refreshToken, err := rtc.RefreshTokenService.CreateRefreshToken(&user, config.CfgToken.RefreshTokenSecret, config.CfgToken.RefreshTokenExpiryHour)
	if err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	refreshTokenResponse := dto.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	result.SuccessResponse[dto.RefreshTokenResponse](c, "Refresh token created successfully", &refreshTokenResponse)
}
