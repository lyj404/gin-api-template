package repositories

import (
	"github.com/lyj404/gin-api-template/domain/entity"
)

type AuditLogRepository interface {
	Create(auditLog *entity.AuditLog) error
	GetByID(id uint) (*entity.AuditLog, error)
	GetByOperator(operatorID uint, page, pageSize int) ([]entity.AuditLog, int64, error)
	GetByTarget(targetType string, targetID uint, page, pageSize int) ([]entity.AuditLog, int64, error)
	GetByTimeRange(startTime, endTime string, page, pageSize int) ([]entity.AuditLog, int64, error)
}
