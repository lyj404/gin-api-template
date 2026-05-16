package service

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/lyj404/gin-api-template/domain/dto"
	"github.com/lyj404/gin-api-template/domain/entity"
	"github.com/lyj404/gin-api-template/domain/repositories"
	"github.com/lyj404/gin-api-template/domain/services"
	"github.com/lyj404/gin-api-template/global"
	"github.com/lyj404/gin-api-template/util"
	"gorm.io/gorm"
)

type userManagementServiceImpl struct {
	userRepo repositories.UserRepository
}

func NewUserManagementService(userRepo repositories.UserRepository) services.UserManagementService {
	return &userManagementServiceImpl{userRepo: userRepo}
}

func (s *userManagementServiceImpl) List(page, pageSize int, keyword string) ([]entity.User, map[uint64][]uint64, map[uint64][]string, int64, error) {
	users, total, err := s.userRepo.List(page, pageSize, keyword)
	if err != nil {
		return nil, nil, nil, 0, err
	}

	roleIDsMap := make(map[uint64][]uint64, len(users))
	roleNamesMap := make(map[uint64][]string, len(users))
	for _, u := range users {
		ids, err := s.userRepo.GetRoleIDsByUserID(u.ID)
		if err != nil {
			return nil, nil, nil, 0, err
		}
		roleIDsMap[u.ID] = ids

		names, err := s.userRepo.GetRoleNamesByUserID(u.ID)
		if err != nil {
			return nil, nil, nil, 0, err
		}
		roleNamesMap[u.ID] = names
	}

	return users, roleIDsMap, roleNamesMap, total, nil
}

func (s *userManagementServiceImpl) GetByID(id uint64) (*entity.User, []uint64, []string, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, nil, nil, err
	}
	roleIDs, err := s.userRepo.GetRoleIDsByUserID(id)
	if err != nil {
		return nil, nil, nil, err
	}
	roleNames, err := s.userRepo.GetRoleNamesByUserID(id)
	if err != nil {
		return nil, nil, nil, err
	}
	return user, roleIDs, roleNames, nil
}

func (s *userManagementServiceImpl) Create(req *dto.CreateUserRequest, operatorID uint64) (*entity.User, error) {
	hashed, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %w", err)
	}

	user := &entity.User{
		Name:     req.Name,
		Email:    req.Email,
		PassWord: hashed,
	}

	err = global.G_DB.Transaction(func(tx *gorm.DB) error {
		var existing entity.User
		if err := tx.Where("email = ?", req.Email).First(&existing).Error; err == nil {
			return errors.New("邮箱已存在")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if err := s.userRepo.Create(tx, user); err != nil {
			return err
		}

		orgUnitID := req.OrgUnitID
		if orgUnitID == 0 {
			var rootOrg entity.OrgUnit
			if err := tx.Where("name = ? AND parent_id IS NULL", "root").First(&rootOrg).Error; err != nil {
				return fmt.Errorf("查找根组织失败: %w", err)
			}
			orgUnitID = rootOrg.ID
		}

		if err := s.userRepo.ReplaceUserRoles(tx, user.ID, orgUnitID, req.RoleIDs); err != nil {
			return err
		}

		afterJSON, _ := json.Marshal(map[string]any{
			"name": user.Name, "email": user.Email, "role_ids": req.RoleIDs,
		})
		return s.audit(tx, operatorID, "create", user.ID, "", string(afterJSON), fmt.Sprintf("创建用户: %s", user.Email))
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userManagementServiceImpl) Update(id uint64, req *dto.UpdateUserRequest, operatorID uint64) (*entity.User, error) {
	old, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	updated := *old
	if req.Name != "" {
		updated.Name = req.Name
	}
	if req.Email != "" {
		updated.Email = req.Email
	}

	err = global.G_DB.Transaction(func(tx *gorm.DB) error {
		if err := s.userRepo.Update(tx, &updated); err != nil {
			return err
		}

		if req.Password != "" {
			hashed, err := util.HashPassword(req.Password)
			if err != nil {
				return err
			}
			if err := s.userRepo.UpdatePassword(tx, id, hashed); err != nil {
				return err
			}
		}

		if req.RoleIDs != nil {
			orgUnitID := req.OrgUnitID
			if orgUnitID == 0 {
				var rootOrg entity.OrgUnit
				if err := tx.Where("name = ? AND parent_id IS NULL", "root").First(&rootOrg).Error; err != nil {
					return fmt.Errorf("查找根组织失败: %w", err)
				}
				orgUnitID = rootOrg.ID
			}
			if err := s.userRepo.ReplaceUserRoles(tx, id, orgUnitID, req.RoleIDs); err != nil {
				return err
			}
		}

		beforeJSON, _ := json.Marshal(map[string]any{"name": old.Name, "email": old.Email})
		afterJSON, _ := json.Marshal(map[string]any{
			"name": updated.Name, "email": updated.Email, "role_ids": req.RoleIDs,
		})
		return s.audit(tx, operatorID, "update", id, string(beforeJSON), string(afterJSON), fmt.Sprintf("更新用户: %s", updated.Email))
	})
	if err != nil {
		return nil, err
	}
	return &updated, nil
}

func (s *userManagementServiceImpl) Delete(id uint64, operatorID uint64) error {
	if id == operatorID {
		return errors.New("不能删除自己")
	}
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	return global.G_DB.Transaction(func(tx *gorm.DB) error {
		if err := s.userRepo.Delete(tx, id); err != nil {
			return err
		}

		beforeJSON, _ := json.Marshal(map[string]any{"name": user.Name, "email": user.Email})
		return s.audit(tx, operatorID, "delete", id, string(beforeJSON), "", fmt.Sprintf("删除用户: %s", user.Email))
	})
}

func (s *userManagementServiceImpl) audit(tx *gorm.DB, operatorID uint64, action string, targetID uint64, before, after, description string) error {
	log := entity.AuditLog{
		OperatorID:   operatorID,
		OperatorName: getOperatorName(tx, operatorID),
		Action:       action,
		TargetType:   "user",
		TargetID:     targetID,
		BeforeData:   before,
		AfterData:    after,
		Description:  description,
	}
	return tx.Create(&log).Error
}
