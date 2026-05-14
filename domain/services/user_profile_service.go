package services

import "github.com/lyj404/gin-api-template/domain/dto"

// ProfileService 个人信息服务接口
type ProfileService interface {
	GetProfile(userID uint) (*dto.ProfileResponse, error)
	UpdateProfile(userID uint, req *dto.UpdateProfileRequest) error
	ChangePassword(userID uint, req *dto.ChangePasswordRequest) error
}
