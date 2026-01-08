package service

import (
	"encoding/json"
	"fmt"

	"github.com/lyj404/gin-api-template/domain/entity"
	"github.com/lyj404/gin-api-template/domain/repositories"
	"github.com/lyj404/gin-api-template/domain/services"
	"github.com/lyj404/gin-api-template/global"
	"gorm.io/gorm"
)

type roleServiceImpl struct {
	roleRepo repositories.RoleRepository
}

func NewRoleService(roleRepo repositories.RoleRepository) services.RoleService {
	return &roleServiceImpl{
		roleRepo: roleRepo,
	}
}

func (s *roleServiceImpl) CreateRole(role *entity.Role, operatorID uint) error {
	return global.G_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(role).Error; err != nil {
			return err
		}

		roleJSON, _ := json.Marshal(role)
		description := fmt.Sprintf("创建角色: %s", role.Name)
		return s.createAuditLog(tx, operatorID, "create", "role", role.ID, "", string(roleJSON), description)
	})
}

func (s *roleServiceImpl) UpdateRole(role *entity.Role, operatorID uint) error {
	return global.G_DB.Transaction(func(tx *gorm.DB) error {
		oldRole, err := s.roleRepo.GetByID(role.ID)
		if err != nil {
			return err
		}

		if err := tx.Save(role).Error; err != nil {
			return err
		}

		oldJSON, _ := json.Marshal(oldRole)
		newJSON, _ := json.Marshal(role)
		description := fmt.Sprintf("更新角色: %s", role.Name)
		return s.createAuditLog(tx, operatorID, "update", "role", role.ID, string(oldJSON), string(newJSON), description)
	})
}

func (s *roleServiceImpl) DeleteRole(id uint, operatorID uint) error {
	return global.G_DB.Transaction(func(tx *gorm.DB) error {
		role, err := s.roleRepo.GetByID(id)
		if err != nil {
			return err
		}

		if role.IsSystem {
			return fmt.Errorf("系统角色不能删除")
		}

		if err := tx.Delete(&entity.Role{}, id).Error; err != nil {
			return err
		}

		roleJSON, _ := json.Marshal(role)
		description := fmt.Sprintf("删除角色: %s", role.Name)
		return s.createAuditLog(tx, operatorID, "delete", "role", id, string(roleJSON), "", description)
	})
}

func (s *roleServiceImpl) GetRoleByID(id uint) (*entity.Role, error) {
	return s.roleRepo.GetByID(id)
}

func (s *roleServiceImpl) GetAllRoles() ([]entity.Role, error) {
	return s.roleRepo.GetAll()
}

func (s *roleServiceImpl) BindResource(roleID, resourceID uint, isWrite bool, operatorID uint) error {
	return global.G_DB.Transaction(func(tx *gorm.DB) error {
		roleResource := entity.RoleResource{
			RoleID:     roleID,
			ResourceID: resourceID,
			IsRead:     true,
			IsWrite:    isWrite,
		}

		if err := tx.Create(&roleResource).Error; err != nil {
			return err
		}

		description := fmt.Sprintf("角色 %d 绑定资源 %d (写权限: %v)", roleID, resourceID, isWrite)
		return s.createAuditLog(tx, operatorID, "bind", "role_resource", roleResource.ID, "", "", description)
	})
}

func (s *roleServiceImpl) UnbindResource(roleID, resourceID uint, operatorID uint) error {
	return global.G_DB.Transaction(func(tx *gorm.DB) error {
		description := fmt.Sprintf("角色 %d 解绑资源 %d", roleID, resourceID)
		if err := tx.Where("role_id = ? AND resource_id = ?", roleID, resourceID).Delete(&entity.RoleResource{}).Error; err != nil {
			return err
		}

		return s.createAuditLog(tx, operatorID, "unbind", "role_resource", 0, "", "", description)
	})
}

func (s *roleServiceImpl) BindOrgScope(roleID, orgUnitID uint, includeDescendants bool, operatorID uint) error {
	return global.G_DB.Transaction(func(tx *gorm.DB) error {
		roleOrgScope := entity.RoleOrgScope{
			RoleID:             roleID,
			OrgUnitID:          orgUnitID,
			IncludeDescendants: includeDescendants,
		}

		if err := tx.Create(&roleOrgScope).Error; err != nil {
			return err
		}

		description := fmt.Sprintf("角色 %d 绑定组织范围 %d (包含子级: %v)", roleID, orgUnitID, includeDescendants)
		return s.createAuditLog(tx, operatorID, "bind", "role_org_scope", roleOrgScope.ID, "", "", description)
	})
}

func (s *roleServiceImpl) UnbindOrgScope(roleID, orgUnitID uint, operatorID uint) error {
	return global.G_DB.Transaction(func(tx *gorm.DB) error {
		description := fmt.Sprintf("角色 %d 解绑组织范围 %d", roleID, orgUnitID)
		if err := tx.Where("role_id = ? AND org_unit_id = ?", roleID, orgUnitID).Delete(&entity.RoleOrgScope{}).Error; err != nil {
			return err
		}

		return s.createAuditLog(tx, operatorID, "unbind", "role_org_scope", 0, "", "", description)
	})
}

func (s *roleServiceImpl) AssignRoleToUser(userID, roleID, orgUnitID uint, operatorID uint) error {
	return global.G_DB.Transaction(func(tx *gorm.DB) error {
		userRole := entity.UserRole{
			UserID:    userID,
			RoleID:    roleID,
			OrgUnitID: orgUnitID,
		}

		if err := tx.Create(&userRole).Error; err != nil {
			return err
		}

		description := fmt.Sprintf("用户 %d 分配角色 %d (组织: %d)", userID, roleID, orgUnitID)
		return s.createAuditLog(tx, operatorID, "assign", "user_role", userRole.ID, "", "", description)
	})
}

func (s *roleServiceImpl) RevokeRoleFromUser(userID, roleID, orgUnitID uint, operatorID uint) error {
	return global.G_DB.Transaction(func(tx *gorm.DB) error {
		description := fmt.Sprintf("用户 %d 撤销角色 %d (组织: %d)", userID, roleID, orgUnitID)
		if err := tx.Where("user_id = ? AND role_id = ? AND org_unit_id = ?", userID, roleID, orgUnitID).Delete(&entity.UserRole{}).Error; err != nil {
			return err
		}

		return s.createAuditLog(tx, operatorID, "revoke", "user_role", 0, "", "", description)
	})
}

func (s *roleServiceImpl) createAuditLog(tx *gorm.DB, operatorID uint, action, targetType string, targetID uint, beforeData, afterData, description string) error {
	auditLog := entity.AuditLog{
		OperatorID:   operatorID,
		OperatorName: getOperatorName(tx, operatorID),
		Action:       action,
		TargetType:   targetType,
		TargetID:     targetID,
		BeforeData:   beforeData,
		AfterData:    afterData,
		Description:  description,
	}
	return tx.Create(&auditLog).Error
}
