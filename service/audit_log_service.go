package service

import (
	"github.com/lyj404/gin-api-template/domain/entity"
	"github.com/lyj404/gin-api-template/domain/repositories"
	"github.com/lyj404/gin-api-template/domain/services"
)

type auditLogServiceImpl struct {
	auditLogRepo repositories.AuditLogRepository
}

func NewAuditLogService(auditLogRepo repositories.AuditLogRepository) services.AuditLogService {
	return &auditLogServiceImpl{
		auditLogRepo: auditLogRepo,
	}
}

func (s *auditLogServiceImpl) GetAuditLogsByOperator(operatorID uint, page, pageSize int) ([]entity.AuditLog, int64, error) {
	return s.auditLogRepo.GetByOperator(operatorID, page, pageSize)
}

func (s *auditLogServiceImpl) GetAuditLogsByTarget(targetType string, targetID uint, page, pageSize int) ([]entity.AuditLog, int64, error) {
	return s.auditLogRepo.GetByTarget(targetType, targetID, page, pageSize)
}

func (s *auditLogServiceImpl) GetAuditLogsByTimeRange(startTime, endTime string, page, pageSize int) ([]entity.AuditLog, int64, error) {
	return s.auditLogRepo.GetByTimeRange(startTime, endTime, page, pageSize)
}
