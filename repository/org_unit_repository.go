package repository

import (
	"fmt"

	"github.com/lyj404/gin-api-template/domain/entity"
	"github.com/lyj404/gin-api-template/domain/repositories"
	"github.com/lyj404/gin-api-template/global"
	"gorm.io/gorm"
)

type orgUnitRepository struct{}

func NewOrgUnitRepository() repositories.OrgUnitRepository {
	return &orgUnitRepository{}
}

func (r *orgUnitRepository) Create(orgUnit *entity.OrgUnit) error {
	// 计算路径和层级
	if orgUnit.ParentID == nil {
		orgUnit.Path = "/" + "0"
		orgUnit.Level = 0
	} else {
		var parent entity.OrgUnit
		if err := global.G_DB.First(&parent, orgUnit.ParentID).Error; err != nil {
			return err
		}
		orgUnit.Path = parent.Path + "/" + "0"
		orgUnit.Level = parent.Level + 1
	}

	return global.G_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(orgUnit).Error; err != nil {
			return err
		}

		// 更新实际ID到路径
		orgUnit.Path = orgUnit.Path[:len(orgUnit.Path)-1] + fmt.Sprintf("%d", orgUnit.ID)
		return tx.Save(orgUnit).Error
	})
}

func (r *orgUnitRepository) Update(orgUnit *entity.OrgUnit) error {
	// 如果修改父节点，需要重新计算路径
	var oldOrg entity.OrgUnit
	if err := global.G_DB.First(&oldOrg, orgUnit.ID).Error; err != nil {
		return err
	}

	if (orgUnit.ParentID == nil && oldOrg.ParentID != nil) || (orgUnit.ParentID != nil && oldOrg.ParentID == nil) || (orgUnit.ParentID != nil && oldOrg.ParentID != nil && *orgUnit.ParentID != *oldOrg.ParentID) {
		return fmt.Errorf("暂不支持修改父节点")
	}

	return global.G_DB.Save(orgUnit).Error
}

func (r *orgUnitRepository) Delete(id uint) error {
	return global.G_DB.Transaction(func(tx *gorm.DB) error {
		// 检查是否有子节点
		var count int64
		if err := tx.Model(&entity.OrgUnit{}).Where("parent_id = ?", id).Count(&count).Error; err != nil {
			return err
		}

		if count > 0 {
			return fmt.Errorf("存在子节点，无法删除")
		}

		return tx.Delete(&entity.OrgUnit{}, id).Error
	})
}

func (r *orgUnitRepository) GetByID(id uint) (*entity.OrgUnit, error) {
	var org entity.OrgUnit
	err := global.G_DB.First(&org, id).Error
	if err != nil {
		return nil, err
	}
	return &org, nil
}

func (r *orgUnitRepository) GetByPath(path string) (*entity.OrgUnit, error) {
	var org entity.OrgUnit
	err := global.G_DB.Where("path = ?", path).First(&org).Error
	if err != nil {
		return nil, err
	}
	return &org, nil
}

func (r *orgUnitRepository) GetAll() ([]entity.OrgUnit, error) {
	var orgs []entity.OrgUnit
	err := global.G_DB.Order("path").Find(&orgs).Error
	return orgs, err
}

func (r *orgUnitRepository) GetChildren(parentID uint) ([]entity.OrgUnit, error) {
	var orgs []entity.OrgUnit
	err := global.G_DB.Where("parent_id = ?", parentID).Order("path").Find(&orgs).Error
	return orgs, err
}
