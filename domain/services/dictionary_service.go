package services

import (
	"context"
	"github.com/lyj404/gin-api-template/domain/entity"
)

// DictionaryService 字典服务接口
type DictionaryService interface {
	CreateDict(ctx context.Context, dict *entity.SysDictionary) error
	UpdateDict(ctx context.Context, dict *entity.SysDictionary) error
	DeleteDict(ctx context.Context, id string) error
	GetDictByID(ctx context.Context, id string) (entity.SysDictionary, error)
	ListDict(ctx context.Context, name, dictType string, status int, page, pageSize int) ([]entity.SysDictionary, int64, error)

	CreateDictDetail(ctx context.Context, detail *entity.SysDictionaryDetail) error
	UpdateDictDetail(ctx context.Context, detail *entity.SysDictionaryDetail) error
	DeleteDictDetail(ctx context.Context, id string) error
	GetDictDetailByID(ctx context.Context, id string) (entity.SysDictionaryDetail, error)
	ListDictDetails(ctx context.Context, dictID string) ([]entity.SysDictionaryDetail, error)

	GetDictInfoByType(ctx context.Context, dictType string) ([]entity.SysDictionaryDetail, error)
}
