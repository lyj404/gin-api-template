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

type resourceServiceImpl struct {
	resourceRepo repositories.ResourceRepository
}

func NewResourceService(resourceRepo repositories.ResourceRepository) services.ResourceService {
	return &resourceServiceImpl{
		resourceRepo: resourceRepo,
	}
}

func (s *resourceServiceImpl) CreateResource(resource *entity.Resource, operatorID uint) error {
	return global.G_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(resource).Error; err != nil {
			return err
		}

		resourceJSON, _ := json.Marshal(resource)
		description := fmt.Sprintf("创建资源: %s", resource.Name)
		return s.createAuditLog(tx, operatorID, "create", "resource", resource.ID, "", string(resourceJSON), description)
	})
}

func (s *resourceServiceImpl) UpdateResource(resource *entity.Resource, operatorID uint) error {
	return global.G_DB.Transaction(func(tx *gorm.DB) error {
		oldResource, err := s.resourceRepo.GetByID(resource.ID)
		if err != nil {
			return err
		}

		if err := tx.Save(resource).Error; err != nil {
			return err
		}

		oldJSON, _ := json.Marshal(oldResource)
		newJSON, _ := json.Marshal(resource)
		description := fmt.Sprintf("更新资源: %s", resource.Name)
		return s.createAuditLog(tx, operatorID, "update", "resource", resource.ID, string(oldJSON), string(newJSON), description)
	})
}

func (s *resourceServiceImpl) DeleteResource(id uint, operatorID uint) error {
	return global.G_DB.Transaction(func(tx *gorm.DB) error {
		resource, err := s.resourceRepo.GetByID(id)
		if err != nil {
			return err
		}

		if err := tx.Delete(&entity.Resource{}, id).Error; err != nil {
			return err
		}

		resourceJSON, _ := json.Marshal(resource)
		description := fmt.Sprintf("删除资源: %s", resource.Name)
		return s.createAuditLog(tx, operatorID, "delete", "resource", id, string(resourceJSON), "", description)
	})
}

func (s *resourceServiceImpl) GetResourceByID(id uint) (*entity.Resource, error) {
	return s.resourceRepo.GetByID(id)
}

func (s *resourceServiceImpl) GetAllResources() ([]entity.Resource, error) {
	return s.resourceRepo.GetAll()
}

func (s *resourceServiceImpl) createAuditLog(tx *gorm.DB, operatorID uint, action, targetType string, targetID uint, beforeData, afterData, description string) error {
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
