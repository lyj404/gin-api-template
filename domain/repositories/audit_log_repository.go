package repositories

import (
	"github.com/lyj404/gin-api-template/domain/entity"
)

type AuditLogRepository interface {
	Create(auditLog *entity.AuditLog) error
	GetByID(id uint64) (*entity.AuditLog, error)
	GetByOperator(operatorID uint64, page, pageSize int, orgIDs []uint64) ([]entity.AuditLog, int64, error)
	GetByTarget(targetType string, targetID uint64, page, pageSize int, orgIDs []uint64) ([]entity.AuditLog, int64, error)
	GetByTimeRange(startTime, endTime string, page, pageSize int, orgIDs []uint64) ([]entity.AuditLog, int64, error)
}
