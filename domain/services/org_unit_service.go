package services

import (
	"github.com/lyj404/gin-api-template/domain/entity"
)

type OrgUnitService interface {
	CreateOrgUnit(orgUnit *entity.OrgUnit, operatorID uint) error
	UpdateOrgUnit(orgUnit *entity.OrgUnit, operatorID uint) error
	DeleteOrgUnit(id uint, operatorID uint) error
	GetOrgUnitByID(id uint) (*entity.OrgUnit, error)
	GetAllOrgUnits() ([]entity.OrgUnit, error)
	GetOrgTree() ([]entity.OrgUnit, error)
}
