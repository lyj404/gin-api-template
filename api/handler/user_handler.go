package handler

import (
	"gin-api-template/config"
	"gin-api-template/domain"
	"gin-api-template/domain/dto"
	"gin-api-template/domain/entity"
	"gin-api-template/domain/result"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	UserService         domain.LoginService
	RefreshTokenUseCase domain.RefreshTokenService
}

func (u *UserHandler) Login(c *gin.Context) {
	var request dto.LoginRequest

	err := c.ShouldBind(&request)
	// 请求传递的参数错误
	if err != nil {
		result.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// 通过邮箱查找用户
	user, err := u.UserService.GetUserByEmail(c, request.Email)
	if err != nil {
		result.ErrorResponse(c, http.StatusNotFound, "User not found with given email")
		return
	}

	// 验证密码
	if bcrypt.CompareHashAndPassword([]byte(user.PassWord), []byte(request.Password)) != nil {
		result.ErrorResponse(c, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// 创建访问token
	accessToken, err := u.RefreshTokenUseCase.CreateAccessToken(&user, config.CfgToken.AccessTokenSecret, config.CfgToken.AccessTokenExpiryHour)
	if err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 创建刷新token
	refreshToken, err := u.RefreshTokenUseCase.CreateRefreshToken(&user, config.CfgToken.RefreshTokenSecret, config.CfgToken.RefreshTokenExpiryHour)
	if err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	loginResponse := dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	// 返回成功响应
	result.SuccessResponse[dto.LoginResponse](c, "Login successful", &loginResponse)
}

func (u *UserHandler) Signup(c *gin.Context) {
	var request dto.SignupRequest
	// 获取请求参数
	err := c.ShouldBind(&request)
	if err != nil {
		result.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// 验证邮箱是否已经存在
	_, err = u.UserService.GetUserByEmail(c, request.Email)
	if err == nil {
		result.ErrorResponse(c, http.StatusConflict, "User already exists with the given email")
		return
	}

	// 加密密码
	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 将加密后的密码赋值给注册信息
	request.Password = string(encryptedPassword)

	user := entity.User{
		Name:     request.Name,
		Email:    request.Email,
		PassWord: request.Password,
	}

	// 将用户数据插入到数据库
	err = u.UserService.Create(c, &user)
	if err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 创建accessToken
	accessToken, err := u.RefreshTokenUseCase.CreateAccessToken(&user, config.CfgToken.AccessTokenSecret, config.CfgToken.AccessTokenExpiryHour)
	if err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 创建refreshToken
	refreshToken, err := u.RefreshTokenUseCase.CreateRefreshToken(&user, config.CfgToken.RefreshTokenSecret, config.CfgToken.RefreshTokenExpiryHour)
	if err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	signupResponse := dto.SignupResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	result.SuccessResponse[dto.SignupResponse](c, "Signup successful", &signupResponse)
}
