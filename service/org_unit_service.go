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

type orgUnitServiceImpl struct {
	orgRepo repositories.OrgUnitRepository
}

func NewOrgUnitService(orgRepo repositories.OrgUnitRepository) services.OrgUnitService {
	return &orgUnitServiceImpl{
		orgRepo: orgRepo,
	}
}

func (s *orgUnitServiceImpl) CreateOrgUnit(orgUnit *entity.OrgUnit, operatorID uint64) error {
	return global.G_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(orgUnit).Error; err != nil {
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

func (s *orgUnitServiceImpl) GetOrgUnitByID(id uint64) (*entity.OrgUnit, error) {
	return s.orgRepo.GetByID(id)
}

func (s *orgUnitServiceImpl) GetAllOrgUnits() ([]entity.OrgUnit, error) {
	return s.orgRepo.GetAll()
}

func (s *orgUnitServiceImpl) GetOrgTree() ([]entity.OrgUnit, error) {
	return s.orgRepo.GetAll()
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
