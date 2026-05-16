package service

import (
	"errors"
	"fmt"

	"github.com/lyj404/gin-api-template/domain/dto"
	"github.com/lyj404/gin-api-template/domain/entity"
	"github.com/lyj404/gin-api-template/domain/services"
	"github.com/lyj404/gin-api-template/global"
	"github.com/lyj404/gin-api-template/util"
	"gorm.io/gorm"
)

type userProfileServiceImpl struct{}

func NewUserProfileService() services.ProfileService {
	return &userProfileServiceImpl{}
}

func (s *userProfileServiceImpl) GetProfile(userID uint64) (*dto.ProfileResponse, error) {
	var user entity.User
	if err := global.G_DB.First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("用户不存在: %w", err)
	}
	return &dto.ProfileResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (s *userProfileServiceImpl) UpdateProfile(userID uint64, req *dto.UpdateProfileRequest) error {
	return global.G_DB.Transaction(func(tx *gorm.DB) error {
		var existing entity.User
		if err := tx.Where("email = ? AND id != ?", req.Email, userID).First(&existing).Error; err == nil {
			return errors.New("邮箱已被占用")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if err := tx.Model(&entity.User{}).Where("id = ?", userID).Updates(map[string]any{
			"name":  req.Name,
			"email": req.Email,
		}).Error; err != nil {
			return err
		}

		return tx.Create(&entity.AuditLog{
			OperatorID:   userID,
			OperatorName: req.Name,
			Action:       "update_profile",
			TargetType:   "user",
			TargetID:     userID,
			Description:  fmt.Sprintf("更新个人信息: %s", req.Email),
		}).Error
	})
}

func (s *userProfileServiceImpl) ChangePassword(userID uint64, req *dto.ChangePasswordRequest) error {
	var user entity.User
	if err := global.G_DB.First(&user, userID).Error; err != nil {
		return fmt.Errorf("用户不存在: %w", err)
	}

	if util.ComparePassword(user.PassWord, req.OldPassword) != nil {
		return errors.New("原密码错误")
	}

	hashed, err := util.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	return global.G_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&entity.User{}).Where("id = ?", userID).Update("password", hashed).Error; err != nil {
			return err
		}
		return tx.Create(&entity.AuditLog{
			OperatorID:   userID,
			OperatorName: user.Name,
			Action:       "change_password",
			TargetType:   "user",
			TargetID:     userID,
			Description:  "修改密码",
		}).Error
	})
}
