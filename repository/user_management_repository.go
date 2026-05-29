package repository

import (
	"github.com/lyj404/gin-api-template/domain/entity"
	"github.com/lyj404/gin-api-template/domain/repositories"
	"github.com/lyj404/gin-api-template/global"
	"gorm.io/gorm"
)

type userManagementRepository struct{}

func NewUserManagementRepository() repositories.UserRepository {
	return &userManagementRepository{}
}

func (r *userManagementRepository) List(page, pageSize int, keyword string, orgIDs []uint64) ([]entity.User, int64, error) {
	var users []entity.User
	var total int64

	// "user" 是 PostgreSQL 保留字，在原始 SQL 中必须用双引号包裹
	query := global.G_DB.Model(&entity.User{}).
		Joins(`JOIN user_role ON user_role.user_id = "user".id AND user_role.deleted_at IS NULL`)
	if orgIDs != nil {
		query = query.Where("user_role.org_unit_id IN ?", orgIDs)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where(`"user".name LIKE ? OR "user".email LIKE ?`, like, like)
	}

	// 使用 WHERE id IN (子查询) 方式计数，避免 DISTINCT+Count/子查询表名引用 在 PostgreSQL 下的兼容问题
	userIDs2 := global.G_DB.Model(&entity.User{}).
		Joins(`JOIN user_role ON user_role.user_id = "user".id AND user_role.deleted_at IS NULL`)
	if orgIDs != nil {
		userIDs2 = userIDs2.Where("user_role.org_unit_id IN ?", orgIDs)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		userIDs2 = userIDs2.Where(`"user".name LIKE ? OR "user".email LIKE ?`, like, like)
	}
	userIDs2 = userIDs2.Select(`DISTINCT "user".id`)
	if err := global.G_DB.Model(&entity.User{}).Where(`"user".id IN (?)`, userIDs2).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Distinct().Order(`"user".id DESC`).Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *userManagementRepository) GetByID(id uint64) (*entity.User, error) {
	var user entity.User
	if err := global.G_DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userManagementRepository) Create(tx *gorm.DB, user *entity.User) error {
	return tx.Create(user).Error
}

func (r *userManagementRepository) Update(tx *gorm.DB, user *entity.User) error {
	return tx.Model(&entity.User{}).Where("id = ?", user.ID).Updates(map[string]any{
		"name":  user.Name,
		"email": user.Email,
	}).Error
}

func (r *userManagementRepository) UpdatePassword(tx *gorm.DB, id uint64, hashed string) error {
	return tx.Model(&entity.User{}).Where("id = ?", id).Update("password", hashed).Error
}

func (r *userManagementRepository) Delete(tx *gorm.DB, id uint64) error {
	if err := tx.Where("user_id = ?", id).Delete(&entity.UserRole{}).Error; err != nil {
		return err
	}
	return tx.Delete(&entity.User{}, id).Error
}

func (r *userManagementRepository) GetRoleIDsByUserID(userID uint64) ([]uint64, error) {
	var roleIDs []uint64
	err := global.G_DB.Model(&entity.UserRole{}).
		Where("user_id = ?", userID).
		Pluck("role_id", &roleIDs).Error
	return roleIDs, err
}

func (r *userManagementRepository) ReplaceUserRoles(tx *gorm.DB, userID, orgUnitID uint64, roleIDs []uint64) error {
	if err := tx.Where("user_id = ?", userID).Delete(&entity.UserRole{}).Error; err != nil {
		return err
	}
	if len(roleIDs) == 0 {
		return nil
	}
	rows := make([]entity.UserRole, 0, len(roleIDs))
	for _, rid := range roleIDs {
		rows = append(rows, entity.UserRole{UserID: userID, RoleID: rid, OrgUnitID: orgUnitID})
	}
	return tx.Create(&rows).Error
}

func (r *userManagementRepository) HasSystemRole(userID uint64) (bool, error) {
	var count int64
	err := global.G_DB.Model(&entity.UserRole{}).
		Joins(`JOIN role ON role.id = user_role.role_id AND role.deleted_at IS NULL`).
		Where("user_role.user_id = ? AND role.is_system = ?", userID, true).
		Count(&count).Error
	return count > 0, err
}

func (r *userManagementRepository) GetRoleNamesByUserID(userID uint64) ([]string, error) {
	var names []string
	err := global.G_DB.Table("user_role").
		Select("role.name").
		Joins("LEFT JOIN role ON role.id = user_role.role_id").
		Where("user_role.user_id = ? AND user_role.deleted_at IS NULL", userID).
		Pluck("role.name", &names).Error
	return names, err
}
