package services

import (
	"github.com/lyj404/gin-api-template/domain/entity"
)

type OrgUnitService interface {
	CreateOrgUnit(orgUnit *entity.OrgUnit, operatorID uint64) error
	UpdateOrgUnit(orgUnit *entity.OrgUnit, operatorID uint64) error
	DeleteOrgUnit(id uint64, operatorID uint64) error
	GetOrgUnitByID(id uint64, userID uint64) (*entity.OrgUnit, error)
	GetAllOrgUnits(userID uint64) ([]entity.OrgUnit, error)
	GetOrgTree(userID uint64) ([]entity.OrgUnit, error)
}
