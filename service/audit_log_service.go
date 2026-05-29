package service

import (
	"github.com/lyj404/gin-api-template/domain/entity"
	"github.com/lyj404/gin-api-template/domain/repositories"
	"github.com/lyj404/gin-api-template/domain/services"
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
	orgIDs := s.getOrgIDs(userID)
	return s.auditLogRepo.GetByOperator(operatorID, page, pageSize, orgIDs)
}

func (s *auditLogServiceImpl) GetAuditLogsByTarget(targetType string, targetID uint64, page, pageSize int, userID uint64) ([]entity.AuditLog, int64, error) {
	orgIDs := s.getOrgIDs(userID)
	return s.auditLogRepo.GetByTarget(targetType, targetID, page, pageSize, orgIDs)
}

func (s *auditLogServiceImpl) GetAuditLogsByTimeRange(startTime, endTime string, page, pageSize int, userID uint64) ([]entity.AuditLog, int64, error) {
	orgIDs := s.getOrgIDs(userID)
	return s.auditLogRepo.GetByTimeRange(startTime, endTime, page, pageSize, orgIDs)
}

func (s *auditLogServiceImpl) getOrgIDs(userID uint64) []uint64 {
	scope, err := s.permSvc.GetUserOrgScope(userID)
	if err != nil {
		return nil
	}
	return CollectOrgIDs(scope)
}
