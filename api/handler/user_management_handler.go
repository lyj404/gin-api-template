package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lyj404/gin-api-template/domain/dto"
	"github.com/lyj404/gin-api-template/domain/result"
	"github.com/lyj404/gin-api-template/domain/services"
)

type UserManagementHandler struct {
	userMgmt services.UserManagementService
}

func NewUserManagementHandler(userMgmt services.UserManagementService) *UserManagementHandler {
	return &UserManagementHandler{userMgmt: userMgmt}
}

// ListUsers 用户列表分页
// @Summary 用户列表
// @Description 获取用户列表（支持分页和关键词搜索）
// @Tags 用户
// @Produce json
// @Param page query int false "页码，默认1"
// @Param page_size query int false "每页数量，默认10，最大100"
// @Param keyword query string false "搜索关键词（搜索用户名或邮箱）"
// @Success 200 {object} result.ResponseResult[dto.PaginationResponse] "获取成功"
// @Failure 400 {object} result.ResponseResult[string] "请求参数错误"
// @Failure 500 {object} result.ResponseResult[string] "服务器内部错误"
// @Router /users [get]
func (h *UserManagementHandler) ListUsers(c *gin.Context) {
	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		result.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	req.SetDefaults()

	users, roleIDsMap, roleNamesMap, total, err := h.userMgmt.List(req.Page, req.PageSize, req.Keyword)
	if err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	responses := make([]dto.UserResponse, len(users))
	for i, u := range users {
		rids := roleIDsMap[u.ID]
		roleIDs := make([]string, len(rids))
		for j, id := range rids {
			roleIDs[j] = strconv.FormatUint(id, 10)
		}
		responses[i] = dto.UserResponse{
			ID:      u.ID,
			Name:    u.Name,
			Email:   u.Email,
			RoleIDs: roleIDs,
			Roles:   roleNamesMap[u.ID],
		}
	}

	result.SuccessResponse(c, "获取用户列表成功", dto.NewPaginationResponse(req.Page, req.PageSize, total, responses))
}

// GetUser 用户详情
// @Summary 用户详情
// @Description 根据ID获取用户详情
// @Tags 用户
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} result.ResponseResult[dto.UserResponse] "获取成功"
// @Failure 400 {object} result.ResponseResult[string] "无效的用户ID"
// @Failure 404 {object} result.ResponseResult[string] "用户不存在"
// @Router /users/{id} [get]
func (h *UserManagementHandler) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		result.ErrorResponse(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	user, roleIDs, roleNames, err := h.userMgmt.GetByID(uint64(id))
	if err != nil {
		result.ErrorResponse(c, http.StatusNotFound, "用户不存在")
		return
	}

	rids := make([]string, len(roleIDs))
	for i, id := range roleIDs {
		rids[i] = strconv.FormatUint(id, 10)
	}
	result.SuccessResponse(c, "获取用户成功", &dto.UserResponse{
		ID:      user.ID,
		Name:    user.Name,
		Email:   user.Email,
		RoleIDs: rids,
		Roles:   roleNames,
	})
}

// CreateUser 创建用户
// @Summary 创建用户
// @Description 创建新用户并分配角色
// @Tags 用户
// @Accept json
// @Produce json
// @Param request body dto.CreateUserRequest true "创建用户请求"
// @Success 200 {object} result.ResponseResult[dto.UserResponse] "创建成功"
// @Failure 400 {object} result.ResponseResult[string] "请求参数错误"
// @Failure 500 {object} result.ResponseResult[string] "服务器内部错误"
// @Router /users [post]
func (h *UserManagementHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	operatorID := c.GetUint64("user_id")
	user, err := h.userMgmt.Create(&req, operatorID)
	if err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.SuccessResponse(c, "用户创建成功", &dto.UserResponse{
		ID:      user.ID,
		Name:    user.Name,
		Email:   user.Email,
		RoleIDs: req.RoleIDs,
	})
}

// UpdateUser 更新用户
// @Summary 更新用户
// @Description 更新用户信息，可同时更新角色绑定
// @Tags 用户
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param request body dto.UpdateUserRequest true "更新用户请求"
// @Success 200 {object} result.ResponseResult[dto.UserResponse] "更新成功"
// @Failure 400 {object} result.ResponseResult[string] "请求参数错误"
// @Failure 404 {object} result.ResponseResult[string] "用户不存在"
// @Failure 500 {object} result.ResponseResult[string] "服务器内部错误"
// @Router /users/{id} [put]
func (h *UserManagementHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		result.ErrorResponse(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	operatorID := c.GetUint64("user_id")
	user, err := h.userMgmt.Update(uint64(id), &req, operatorID)
	if err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.SuccessResponse(c, "用户更新成功", &dto.UserResponse{
		ID:      user.ID,
		Name:    user.Name,
		Email:   user.Email,
		RoleIDs: req.RoleIDs,
	})
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Description 根据ID删除用户
// @Tags 用户
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} result.ResponseResult[string] "删除成功"
// @Failure 400 {object} result.ResponseResult[string] "无效的用户ID"
// @Failure 500 {object} result.ResponseResult[string] "服务器内部错误"
// @Router /users/{id} [delete]
func (h *UserManagementHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		result.ErrorResponse(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	operatorID := c.GetUint64("user_id")
	if uint64(id) == operatorID {
		result.ErrorResponse(c, http.StatusBadRequest, "不能删除自己")
		return
	}
	if err := h.userMgmt.Delete(uint64(id), operatorID); err != nil {
		result.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.SimpleSuccessResponse(c, "用户删除成功")
}
