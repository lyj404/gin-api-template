package handler

import (
	"gin-api-template/config"
	"gin-api-template/domain"
	"gin-api-template/domain/dto"
	"gin-api-template/domain/result"
	"log"
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
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	// 通过邮箱查找用户
	user, err := u.UserService.GetUserByEmail(c, request.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, result.ErrorResponse{
			Message: "User not found with given email",
		})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PassWord), []byte(request.Password)) != nil {
		c.JSON(http.StatusUnauthorized, result.ErrorResponse{
			Message: "Invalid credentials",
		})
		return
	}

	accessToken, err := u.RefreshTokenUseCase.CreateAccessToken(&user, config.AccessTokenSecret, config.AccessTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	refreshToken, err := u.RefreshTokenUseCase.CreateRefreshToken(&user, config.RefreshTokenSecret, config.RefreshTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	loginResponse := dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	c.JSON(http.StatusOK, loginResponse)
}

func (u *UserHandler) Signup(c *gin.Context) {
	var request dto.SignupRequest
	log.Println("Signup request received")
	// 获取注册信息
	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorResponse{Message: err.Error()})
		return
	}

	// 验证邮箱是否已经存在
	_, err = u.UserService.GetUserByEmail(c, request.Email)
	if err == nil {
		c.JSON(http.StatusConflict, result.ErrorResponse{Message: "User already exists with the given email"})
		return
	}

	// 加密密码
	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorResponse{Message: err.Error()})
		return
	}

	// 将加密后的密码赋值给注册信息
	request.Password = string(encryptedPassword)

	user := domain.User{
		Name:     request.Name,
		Email:    request.Email,
		PassWord: request.Password,
	}

	// 将用户数据插入到数据库
	err = u.UserService.Create(c, &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorResponse{Message: err.Error()})
		return
	}

	// 创建accessToken
	accessToken, err := u.RefreshTokenUseCase.CreateAccessToken(&user, config.AccessTokenSecret, config.AccessTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorResponse{Message: err.Error()})
		return
	}

	// 创建refreshToken
	refreshToken, err := u.RefreshTokenUseCase.CreateRefreshToken(&user, config.RefreshTokenSecret, config.RefreshTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorResponse{Message: err.Error()})
		return
	}

	signupResponse := dto.SignupResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	c.JSON(http.StatusOK, signupResponse)
}
