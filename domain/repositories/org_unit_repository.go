package repositories

import (
	"github.com/lyj404/gin-api-template/domain/entity"
)

type OrgUnitRepository interface {
	Create(orgUnit *entity.OrgUnit) error
	Update(orgUnit *entity.OrgUnit) error
	Delete(id uint) error
	GetByID(id uint) (*entity.OrgUnit, error)
	GetByPath(path string) (*entity.OrgUnit, error)
	GetAll() ([]entity.OrgUnit, error)
	GetChildren(parentID uint) ([]entity.OrgUnit, error)
}
