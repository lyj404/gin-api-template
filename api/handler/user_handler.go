package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/lyj404/gin-api-template/config"
	"github.com/lyj404/gin-api-template/domain"
	"github.com/lyj404/gin-api-template/domain/dto"
	"github.com/lyj404/gin-api-template/domain/entity"
	"github.com/lyj404/gin-api-template/domain/result"
	"github.com/lyj404/gin-api-template/pkg/lib/captcha"
	"github.com/lyj404/gin-api-template/util"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserService         domain.LoginService
	RefreshTokenUseCase domain.RefreshTokenService
}

// @Summary 用户登录
// @Description 处理用户登录请求，验证凭据并返回访问token和刷新token
// @Tags user
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "登录请求参数"
// @Success 200 {object} result.ResponseResult[dto.LoginResponse] "登录成功响应"
// @Failure 400 {object} result.ResponseResult[string] "请求参数错误"
// @Failure 404 {object} result.ResponseResult[string] "用户未找到"
// @Failure 401 {object} result.ResponseResult[string] "凭据无效"
// @Failure 500 {object} result.ResponseResult[string] "服务器内部错误"
// @Router /login [post]
func (u *UserHandler) Login(c *gin.Context) {
	var request dto.LoginRequest

	err := c.ShouldBind(&request)
	// 请求传递的参数错误
	if err != nil {
		result.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// 验证验证码
	session, err := captcha.Store.Get(c.Request, "captcha-session")
	if err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, "Session Error")
	}
	captchaData, ok := session.Values[captcha.CaptchaSessionKey].(map[string]interface{})
	if !ok {
		result.ErrorResponse(c, http.StatusBadRequest, "Verification code has expired, please refresh")
		return
	}

	// 检查时间有效性
	generateTime, ok := captchaData["time"].(int64)
	if !ok || time.Now().Unix()-generateTime > 300 {
		result.ErrorResponse(c, http.StatusBadRequest, "Verification code has expired, please refresh")
		return
	}

	storedAnswer, ok := captchaData["answer"].(int)
	if !ok {
		result.ErrorResponse(c, http.StatusInternalServerError, "Verification code data error")
		return
	}

	userAnswer, err := strconv.Atoi(strings.TrimSpace(request.Captcha))
	if err != nil || userAnswer != storedAnswer {
		result.ErrorResponse(c, http.StatusBadRequest, "Incorrect verification code")
		return
	}

	// 通过邮箱查找用户
	user, err := u.UserService.GetUserByEmail(c, request.Email)
	if err != nil {
		result.ErrorResponse(c, http.StatusNotFound, "User not found with given email")
		return
	}

	// 验证密码
	if util.ComparePassword(user.PassWord, request.Password) != nil {
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
	result.SuccessResponse(c, "Login successful", &loginResponse)
}

// @Summary 用户注册
// @Description 处理新用户注册请求，创建用户账户并返回访问token和刷新token
// @Tags user
// @Accept json
// @Produce json
// @Param request body dto.SignupRequest true "注册请求参数"
// @Success 200 {object} result.ResponseResult[dto.SignupResponse] "注册成功响应"
// @Failure 400 {object} result.ResponseResult[string] "请求参数错误"
// @Failure 409 {object} result.ResponseResult[string] "邮箱已存在"
// @Failure 500 {object} result.ResponseResult[string] "服务器内部错误"
// @Router /signup [post]
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
	encryptedPassword, err := util.HashPassword(request.Password)
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

	result.SuccessResponse(c, "Signup successful", &signupResponse)
}

// @Summary 验证码
// @Description 生成验证码防止反复登录
// @Tags user
// @Accept json
// @Produce json
// @Param request body dto.SignupRequest true "注册请求参数"
// @Success 200 {object} result.ResponseResult[captcha.CaptchaReponse] "验证码成功响应"
// @Failure 500 {object} result.ResponseResult[string] "服务器内部错误"
// @Router /signup [post]
func (u *UserHandler) GenerateMathCaptcha(c *gin.Context) {
	// 获取 session
	session, err := captcha.Store.Get(c.Request, "captcha-session")
	if err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, "Session Error")
	}

	// 生成数学问题
	problem := captcha.GenerateMathProblem()
	// 生成验证码图片
	imageURL, err := captcha.GenerateCaptchaImage(problem.Question)
	if err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate verification code image")
	}
	// 存储问题到 session
	session.Values[captcha.CaptchaSessionKey] = map[string]interface{}{
		"answer":   problem.Answer,
		"question": problem.Question,
		"type":     int(problem.Type),
		"time":     time.Now().Unix(),
	}
	// 保存 session
	if err := session.Save(c.Request, c.Writer); err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, "Failed to save session")
		return
	}
	response := captcha.CaptchaReponse{
		ImageUrl:   imageURL,
		Question:   problem.Question,
		ExpireTime: 300,
	}
	result.SuccessResponse(c, "success", &response)
}
