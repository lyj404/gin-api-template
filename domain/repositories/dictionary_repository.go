package repositories

import (
	"context"
	"github.com/lyj404/gin-api-template/domain/entity"
)

// DictionaryRepo 字典仓储接口
type DictionaryRepo interface {
	CreateDict(ctx context.Context, dict *entity.SysDictionary) error
	UpdateDict(ctx context.Context, dict *entity.SysDictionary) error
	DeleteDict(ctx context.Context, id string) error
	GetDictByID(ctx context.Context, id string) (entity.SysDictionary, error)
	GetDictByType(ctx context.Context, dictType string) (entity.SysDictionary, error)
	ListDict(ctx context.Context, name, dictType string) ([]entity.SysDictionary, error)

	CreateDictDetail(ctx context.Context, detail *entity.SysDictionaryDetail) error
	UpdateDictDetail(ctx context.Context, detail *entity.SysDictionaryDetail) error
	DeleteDictDetail(ctx context.Context, id string) error
	GetDictDetailByID(ctx context.Context, id string) (entity.SysDictionaryDetail, error)
	ListDictDetails(ctx context.Context, dictID string) ([]entity.SysDictionaryDetail, error)
}
