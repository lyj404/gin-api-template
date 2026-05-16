package repository

import (
	"context"

	"github.com/lyj404/gin-api-template/domain/entity"
	"github.com/lyj404/gin-api-template/domain/repositories"

	"gorm.io/gorm"
)

type dictionaryRepo struct {
	db *gorm.DB
}

func NewDictionaryRepo(db *gorm.DB) repositories.DictionaryRepo {
	return &dictionaryRepo{db: db}
}

func (r *dictionaryRepo) CreateDict(ctx context.Context, dict *entity.SysDictionary) error {
	return r.db.WithContext(ctx).Create(dict).Error
}

func (r *dictionaryRepo) UpdateDict(ctx context.Context, dict *entity.SysDictionary) error {
	return r.db.WithContext(ctx).Save(dict).Error
}

func (r *dictionaryRepo) DeleteDict(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("dict_id = ?", id).Delete(&entity.SysDictionaryDetail{}).Error; err != nil {
			return err
		}
		return tx.Delete(&entity.SysDictionary{}, id).Error
	})
}

func (r *dictionaryRepo) GetDictByID(ctx context.Context, id string) (entity.SysDictionary, error) {
	var dict entity.SysDictionary
	err := r.db.WithContext(ctx).Preload("Details").First(&dict, id).Error
	return dict, err
}

func (r *dictionaryRepo) GetDictByType(ctx context.Context, dictType string) (entity.SysDictionary, error) {
	var dict entity.SysDictionary
	err := r.db.WithContext(ctx).Preload("Details", func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", 1).Order("sort asc")
	}).Where("type = ? AND status = ?", dictType, 1).First(&dict).Error
	return dict, err
}

func (r *dictionaryRepo) ListDict(ctx context.Context, name, dictType string) ([]entity.SysDictionary, error) {
	var dicts []entity.SysDictionary
	db := r.db.WithContext(ctx)
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	if dictType != "" {
		db = db.Where("type LIKE ?", "%"+dictType+"%")
	}
	err := db.Find(&dicts).Error
	return dicts, err
}

func (r *dictionaryRepo) CreateDictDetail(ctx context.Context, detail *entity.SysDictionaryDetail) error {
	return r.db.WithContext(ctx).Create(detail).Error
}

func (r *dictionaryRepo) UpdateDictDetail(ctx context.Context, detail *entity.SysDictionaryDetail) error {
	return r.db.WithContext(ctx).Save(detail).Error
}

func (r *dictionaryRepo) DeleteDictDetail(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entity.SysDictionaryDetail{}, id).Error
}

func (r *dictionaryRepo) GetDictDetailByID(ctx context.Context, id string) (entity.SysDictionaryDetail, error) {
	var detail entity.SysDictionaryDetail
	err := r.db.WithContext(ctx).First(&detail, id).Error
	return detail, err
}

func (r *dictionaryRepo) ListDictDetails(ctx context.Context, dictID string) ([]entity.SysDictionaryDetail, error) {
	var details []entity.SysDictionaryDetail
	err := r.db.WithContext(ctx).Where("dict_id = ?", dictID).Order("sort asc").Find(&details).Error
	return details, err
}
