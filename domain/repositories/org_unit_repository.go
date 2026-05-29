package repositories

import (
	"github.com/lyj404/gin-api-template/domain/entity"
	"gorm.io/gorm"
)

type OrgUnitRepository interface {
	Create(tx *gorm.DB, orgUnit *entity.OrgUnit) error
	Update(orgUnit *entity.OrgUnit) error
	Delete(id uint64) error
	GetByID(id uint64) (*entity.OrgUnit, error)
	GetByPath(path string) (*entity.OrgUnit, error)
	GetAll() ([]entity.OrgUnit, error)
	GetChildren(parentID uint64) ([]entity.OrgUnit, error)
}
