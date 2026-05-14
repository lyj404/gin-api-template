package services

import (
	"github.com/lyj404/gin-api-template/domain/entity"
)

type AuditLogService interface {
	Create(auditLog *entity.AuditLog) error
	GetAuditLogsByOperator(operatorID uint, page, pageSize int) ([]entity.AuditLog, int64, error)
	GetAuditLogsByTarget(targetType string, targetID uint, page, pageSize int) ([]entity.AuditLog, int64, error)
	GetAuditLogsByTimeRange(startTime, endTime string, page, pageSize int) ([]entity.AuditLog, int64, error)
}
