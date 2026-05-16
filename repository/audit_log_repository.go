package repository

import (
	"github.com/lyj404/gin-api-template/domain/entity"
	"github.com/lyj404/gin-api-template/domain/repositories"
	"github.com/lyj404/gin-api-template/global"
)

type auditLogRepository struct{}

func NewAuditLogRepository() repositories.AuditLogRepository {
	return &auditLogRepository{}
}

func (r *auditLogRepository) Create(auditLog *entity.AuditLog) error {
	return global.G_DB.Create(auditLog).Error
}

func (r *auditLogRepository) GetByID(id uint64) (*entity.AuditLog, error) {
	var auditLog entity.AuditLog
	err := global.G_DB.First(&auditLog, id).Error
	if err != nil {
		return nil, err
	}
	return &auditLog, nil
}

func (r *auditLogRepository) GetByOperator(operatorID uint64, page, pageSize int) ([]entity.AuditLog, int64, error) {
	var auditLogs []entity.AuditLog
	var total int64

	offset := (page - 1) * pageSize
	err := global.G_DB.Model(&entity.AuditLog{}).
		Where("operator_id = ?", operatorID).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&auditLogs).Error

	global.G_DB.Model(&entity.AuditLog{}).Where("operator_id = ?", operatorID).Count(&total)

	return auditLogs, total, err
}

func (r *auditLogRepository) GetByTarget(targetType string, targetID uint64, page, pageSize int) ([]entity.AuditLog, int64, error) {
	var auditLogs []entity.AuditLog
	var total int64

	offset := (page - 1) * pageSize
	err := global.G_DB.Model(&entity.AuditLog{}).
		Where("target_type = ? AND target_id = ?", targetType, targetID).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&auditLogs).Error

	global.G_DB.Model(&entity.AuditLog{}).
		Where("target_type = ? AND target_id = ?", targetType, targetID).
		Count(&total)

	return auditLogs, total, err
}

func (r *auditLogRepository) GetByTimeRange(startTime, endTime string, page, pageSize int) ([]entity.AuditLog, int64, error) {
	var auditLogs []entity.AuditLog
	var total int64

	offset := (page - 1) * pageSize
	query := global.G_DB.Model(&entity.AuditLog{})

	if startTime != "" {
		query = query.Where("created_at >= ?", startTime)
	}
	if endTime != "" {
		query = query.Where("created_at <= ?", endTime)
	}

	err := query.
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&auditLogs).Error

	query.Count(&total)

	return auditLogs, total, err
}
