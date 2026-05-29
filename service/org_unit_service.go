package service

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/lyj404/gin-api-template/domain/entity"
	"github.com/lyj404/gin-api-template/domain/repositories"
	"github.com/lyj404/gin-api-template/domain/services"
	"github.com/lyj404/gin-api-template/global"
	"gorm.io/gorm"
)

type orgUnitServiceImpl struct {
	orgRepo repositories.OrgUnitRepository
	permSvc services.PermissionService
}

func NewOrgUnitService(orgRepo repositories.OrgUnitRepository, permSvc services.PermissionService) services.OrgUnitService {
	return &orgUnitServiceImpl{
		orgRepo: orgRepo,
		permSvc: permSvc,
	}
}

func (s *orgUnitServiceImpl) CreateOrgUnit(orgUnit *entity.OrgUnit, operatorID uint64) error {
	return global.G_DB.Transaction(func(tx *gorm.DB) error {
		// 通过 repository 创建，正确计算 Path 和 Level
		if err := s.orgRepo.Create(tx, orgUnit); err != nil {
			return err
		}

		orgJSON, _ := json.Marshal(orgUnit)
		description := fmt.Sprintf("创建组织节点: %s", orgUnit.Name)
		return s.createAuditLog(tx, operatorID, "create", "org_unit", orgUnit.ID, "", string(orgJSON), description)
	})
}

func (s *orgUnitServiceImpl) UpdateOrgUnit(orgUnit *entity.OrgUnit, operatorID uint64) error {
	return global.G_DB.Transaction(func(tx *gorm.DB) error {
		oldOrg, err := s.orgRepo.GetByID(orgUnit.ID)
		if err != nil {
			return err
		}

		if err := tx.Save(orgUnit).Error; err != nil {
			return err
		}

		oldJSON, _ := json.Marshal(oldOrg)
		newJSON, _ := json.Marshal(orgUnit)
		description := fmt.Sprintf("更新组织节点: %s", orgUnit.Name)
		return s.createAuditLog(tx, operatorID, "update", "org_unit", orgUnit.ID, string(oldJSON), string(newJSON), description)
	})
}

func (s *orgUnitServiceImpl) DeleteOrgUnit(id uint64, operatorID uint64) error {
	return global.G_DB.Transaction(func(tx *gorm.DB) error {
		org, err := s.orgRepo.GetByID(id)
		if err != nil {
			return err
		}

		if err := tx.Delete(&entity.OrgUnit{}, id).Error; err != nil {
			return err
		}

		orgJSON, _ := json.Marshal(org)
		description := fmt.Sprintf("删除组织节点: %s", org.Name)
		return s.createAuditLog(tx, operatorID, "delete", "org_unit", id, string(orgJSON), "", description)
	})
}

func (s *orgUnitServiceImpl) GetOrgUnitByID(id uint64, userID uint64) (*entity.OrgUnit, error) {
	org, err := s.orgRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	scope, err := s.permSvc.GetUserOrgScope(userID)
	if err != nil || !s.orgInScope(org, scope) {
		return nil, fmt.Errorf("组织不存在")
	}
	return org, nil
}

func (s *orgUnitServiceImpl) GetAllOrgUnits(userID uint64) ([]entity.OrgUnit, error) {
	all, err := s.orgRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return s.filterByScope(all, userID), nil
}

func (s *orgUnitServiceImpl) GetOrgTree(userID uint64) ([]entity.OrgUnit, error) {
	all, err := s.orgRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return s.filterByScope(all, userID), nil
}

func (s *orgUnitServiceImpl) filterByScope(orgs []entity.OrgUnit, userID uint64) []entity.OrgUnit {
	scope, err := s.permSvc.GetUserOrgScope(userID)
	if err != nil || len(scope) == 0 {
		return nil
	}

	result := make([]entity.OrgUnit, 0, len(orgs))
	for _, org := range orgs {
		if s.orgInScope(&org, scope) {
			result = append(result, org)
		}
	}
	return result
}

func (s *orgUnitServiceImpl) orgInScope(org *entity.OrgUnit, scope []services.OrgScopeInfo) bool {
	for _, sc := range scope {
		if sc.IncludeDescendants {
			if org.ID == sc.OrgUnitID || strings.HasPrefix(org.Path, sc.Path+"/") {
				return true
			}
		} else if org.ID == sc.OrgUnitID {
			return true
		}
	}
	return false
}

func (s *orgUnitServiceImpl) createAuditLog(tx *gorm.DB, operatorID uint64, action, targetType string, targetID uint64, beforeData, afterData, description string) error {
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
