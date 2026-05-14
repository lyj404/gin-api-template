package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lyj404/gin-api-template/domain/dto"
	"github.com/lyj404/gin-api-template/domain/result"
	"github.com/lyj404/gin-api-template/domain/services"
)

type UserProfileHandler struct {
	profileService services.ProfileService
}

func NewUserProfileHandler(profileService services.ProfileService) *UserProfileHandler {
	return &UserProfileHandler{profileService: profileService}
}

// GetProfile 获取个人信息
// @Summary 获取个人信息
// @Description 获取当前登录用户的个人信息
// @Tags 用户
// @Produce json
// @Success 200 {object} result.ResponseResult[dto.ProfileResponse] "获取成功"
// @Failure 401 {object} result.ResponseResult[string] "未授权"
// @Failure 500 {object} result.ResponseResult[string] "服务器内部错误"
// @Router /user/profile [get]
func (h *UserProfileHandler) GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	profile, err := h.profileService.GetProfile(userID)
	if err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	result.SuccessResponse(c, "获取个人信息成功", profile)
}

// UpdateProfile 更新个人信息
// @Summary 更新个人信息
// @Description 更新当前登录用户的姓名和邮箱
// @Tags 用户
// @Accept json
// @Produce json
// @Param request body dto.UpdateProfileRequest true "更新个人信息请求"
// @Success 200 {object} result.ResponseResult[string] "更新成功"
// @Failure 400 {object} result.ResponseResult[string] "请求参数错误"
// @Failure 401 {object} result.ResponseResult[string] "未授权"
// @Failure 500 {object} result.ResponseResult[string] "服务器内部错误"
// @Router /user/profile [put]
func (h *UserProfileHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.profileService.UpdateProfile(userID, &req); err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.SimpleSuccessResponse(c, "个人信息更新成功")
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 修改当前登录用户的密码
// @Tags 用户
// @Accept json
// @Produce json
// @Param request body dto.ChangePasswordRequest true "修改密码请求"
// @Success 200 {object} result.ResponseResult[string] "修改成功"
// @Failure 400 {object} result.ResponseResult[string] "请求参数错误或原密码错误"
// @Failure 401 {object} result.ResponseResult[string] "未授权"
// @Failure 500 {object} result.ResponseResult[string] "服务器内部错误"
// @Router /user/password [put]
func (h *UserProfileHandler) ChangePassword(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.profileService.ChangePassword(userID, &req); err != nil {
		result.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	result.SimpleSuccessResponse(c, "密码修改成功")
}
