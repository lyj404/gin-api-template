package repository

import (
	"github.com/lyj404/gin-api-template/domain/entity"
	"github.com/lyj404/gin-api-template/domain/repositories"
	"github.com/lyj404/gin-api-template/global"
	"gorm.io/gorm"
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

func (r *auditLogRepository) GetByOperator(operatorID uint64, page, pageSize int, orgIDs []uint64) ([]entity.AuditLog, int64, error) {
	var auditLogs []entity.AuditLog
	var total int64

	offset := (page - 1) * pageSize
	query := r.scopedQuery(orgIDs).
		Order("audit_log.created_at DESC").
		Limit(pageSize).
		Offset(offset)

	if operatorID != 0 {
		query = query.Where("audit_log.operator_id = ?", operatorID)
	}

	err := query.Find(&auditLogs).Error
	if err != nil {
		return nil, 0, err
	}

	countQuery := r.scopedQuery(orgIDs)
	if operatorID != 0 {
		countQuery = countQuery.Where("audit_log.operator_id = ?", operatorID)
	}
	countQuery.Count(&total)

	return auditLogs, total, nil
}

func (r *auditLogRepository) GetByTarget(targetType string, targetID uint64, page, pageSize int, orgIDs []uint64) ([]entity.AuditLog, int64, error) {
	var auditLogs []entity.AuditLog
	var total int64

	offset := (page - 1) * pageSize
	err := r.scopedQuery(orgIDs).
		Where("audit_log.target_type = ? AND audit_log.target_id = ?", targetType, targetID).
		Order("audit_log.created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&auditLogs).Error
	if err != nil {
		return nil, 0, err
	}

	r.scopedQuery(orgIDs).
		Where("audit_log.target_type = ? AND audit_log.target_id = ?", targetType, targetID).
		Count(&total)

	return auditLogs, total, nil
}

func (r *auditLogRepository) GetByTimeRange(startTime, endTime string, page, pageSize int, orgIDs []uint64) ([]entity.AuditLog, int64, error) {
	var auditLogs []entity.AuditLog
	var total int64

	offset := (page - 1) * pageSize
	query := r.scopedQuery(orgIDs)

	if startTime != "" {
		query = query.Where("audit_log.created_at >= ?", startTime)
	}
	if endTime != "" {
		query = query.Where("audit_log.created_at <= ?", endTime)
	}

	err := query.
		Order("audit_log.created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&auditLogs).Error
	if err != nil {
		return nil, 0, err
	}

	countQuery := r.scopedQuery(orgIDs)
	if startTime != "" {
		countQuery = countQuery.Where("audit_log.created_at >= ?", startTime)
	}
	if endTime != "" {
		countQuery = countQuery.Where("audit_log.created_at <= ?", endTime)
	}
	countQuery.Count(&total)

	return auditLogs, total, nil
}

func (r *auditLogRepository) scopedQuery(orgIDs []uint64) *gorm.DB {
	// 使用 WHERE EXISTS 子查询替代 JOIN，避免操作者有多角色时产生重复行
	return global.G_DB.Model(&entity.AuditLog{}).
		Where(`EXISTS (SELECT 1 FROM user_role WHERE user_role.user_id = audit_log.operator_id AND user_role.org_unit_id IN ? AND user_role.deleted_at IS NULL)`, orgIDs)
}
