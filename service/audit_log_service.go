package service

import (
	"github.com/lyj404/gin-api-template/domain/entity"
	"github.com/lyj404/gin-api-template/domain/repositories"
	"github.com/lyj404/gin-api-template/domain/services"
	"github.com/lyj404/gin-api-template/global"
)

type auditLogServiceImpl struct {
	auditLogRepo repositories.AuditLogRepository
	permSvc      services.PermissionService
}

func NewAuditLogService(auditLogRepo repositories.AuditLogRepository, permSvc services.PermissionService) services.AuditLogService {
	return &auditLogServiceImpl{
		auditLogRepo: auditLogRepo,
		permSvc:      permSvc,
	}
}

func (s *auditLogServiceImpl) Create(auditLog *entity.AuditLog) error {
	return s.auditLogRepo.Create(auditLog)
}

func (s *auditLogServiceImpl) GetAuditLogsByOperator(operatorID uint64, page, pageSize int, userID uint64) ([]entity.AuditLog, int64, error) {
	orgIDs, err := s.getOrgIDs(userID)
	if err != nil {
		return nil, 0, err
	}
	return s.auditLogRepo.GetByOperator(operatorID, page, pageSize, orgIDs)
}

func (s *auditLogServiceImpl) GetAuditLogsByTarget(targetType string, targetID uint64, page, pageSize int, userID uint64) ([]entity.AuditLog, int64, error) {
	orgIDs, err := s.getOrgIDs(userID)
	if err != nil {
		return nil, 0, err
	}
	return s.auditLogRepo.GetByTarget(targetType, targetID, page, pageSize, orgIDs)
}

func (s *auditLogServiceImpl) GetAuditLogsByTimeRange(startTime, endTime string, page, pageSize int, userID uint64) ([]entity.AuditLog, int64, error) {
	orgIDs, err := s.getOrgIDs(userID)
	if err != nil {
		return nil, 0, err
	}
	return s.auditLogRepo.GetByTimeRange(startTime, endTime, page, pageSize, orgIDs)
}

func (s *auditLogServiceImpl) getOrgIDs(userID uint64) ([]uint64, error) {
	isSuper, err := s.permSvc.HasSystemRole(userID)
	if err != nil {
		return nil, err
	}
	if isSuper {
		return nil, nil
	}
	scope, err := s.permSvc.GetUserOrgScope(userID)
	if err != nil {
		return nil, err
	}
	orgIDs := CollectOrgIDs(scope)
	if len(orgIDs) == 0 {
		var userRole entity.UserRole
		if err := global.G_DB.Where("user_id = ?", userID).First(&userRole).Error; err == nil {
			orgIDs = []uint64{userRole.OrgUnitID}
		}
	}
	return orgIDs, nil
}
