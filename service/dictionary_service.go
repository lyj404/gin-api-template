package service

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/lyj404/gin-api-template/domain/entity"
	"github.com/lyj404/gin-api-template/domain/repositories"
	"github.com/lyj404/gin-api-template/domain/services"
	"github.com/lyj404/gin-api-template/global"
)

type dictionaryService struct {
	repo repositories.DictionaryRepo
}

func NewDictionaryService(repo repositories.DictionaryRepo) services.DictionaryService {
	return &dictionaryService{
		repo: repo,
	}
}

const (
	dictCachePrefix = "dict:"
	dictCacheExpire = 24 * time.Hour
)

func (s *dictionaryService) CreateDict(ctx context.Context, dict *entity.SysDictionary) error {
	return s.repo.CreateDict(ctx, dict)
}

func (s *dictionaryService) UpdateDict(ctx context.Context, dict *entity.SysDictionary) error {
	err := s.repo.UpdateDict(ctx, dict)
	if err == nil {
		s.clearCache(ctx, dict.Type)
	}
	return err
}

func (s *dictionaryService) DeleteDict(ctx context.Context, id string) error {
	dict, err := s.repo.GetDictByID(ctx, id)
	if err != nil {
		return err
	}
	err = s.repo.DeleteDict(ctx, id)
	if err == nil {
		s.clearCache(ctx, dict.Type)
	}
	return err
}

func (s *dictionaryService) GetDictByID(ctx context.Context, id string) (entity.SysDictionary, error) {
	return s.repo.GetDictByID(ctx, id)
}

func (s *dictionaryService) ListDict(ctx context.Context, name, dictType string) ([]entity.SysDictionary, error) {
	return s.repo.ListDict(ctx, name, dictType)
}

func (s *dictionaryService) CreateDictDetail(ctx context.Context, detail *entity.SysDictionaryDetail) error {
	err := s.repo.CreateDictDetail(ctx, detail)
	if err == nil {
		dict, _ := s.repo.GetDictByID(ctx, strconv.FormatUint(detail.DictID, 10))
		s.clearCache(ctx, dict.Type)
	}
	return err
}

func (s *dictionaryService) UpdateDictDetail(ctx context.Context, detail *entity.SysDictionaryDetail) error {
	err := s.repo.UpdateDictDetail(ctx, detail)
	if err == nil {
		dict, _ := s.repo.GetDictByID(ctx, strconv.FormatUint(detail.DictID, 10))
		s.clearCache(ctx, dict.Type)
	}
	return err
}

func (s *dictionaryService) DeleteDictDetail(ctx context.Context, id string) error {
	detail, err := s.repo.GetDictDetailByID(ctx, id)
	if err != nil {
		return err
	}
	err = s.repo.DeleteDictDetail(ctx, id)
	if err == nil {
		dict, _ := s.repo.GetDictByID(ctx, strconv.FormatUint(detail.DictID, 10))
		s.clearCache(ctx, dict.Type)
	}
	return err
}

func (s *dictionaryService) ListDictDetails(ctx context.Context, dictID string) ([]entity.SysDictionaryDetail, error) {
	return s.repo.ListDictDetails(ctx, dictID)
}

func (s *dictionaryService) GetDictDetailByID(ctx context.Context, id string) (entity.SysDictionaryDetail, error) {
	return s.repo.GetDictDetailByID(ctx, id)
}

func (s *dictionaryService) GetDictInfoByType(ctx context.Context, dictType string) ([]entity.SysDictionaryDetail, error) {
	// 尝试从缓存获取
	if global.G_REDIS != nil {
		cacheKey := dictCachePrefix + dictType
		val, err := global.G_REDIS.Get(ctx, cacheKey).Result()
		if err == nil {
			var details []entity.SysDictionaryDetail
			if err := json.Unmarshal([]byte(val), &details); err == nil {
				return details, nil
			}
		}
	}

	// 数据库查询
	dict, err := s.repo.GetDictByType(ctx, dictType)
	if err != nil {
		return nil, err
	}

	// 写入缓存
	if global.G_REDIS != nil {
		cacheKey := dictCachePrefix + dictType
		data, _ := json.Marshal(dict.Details)
		global.G_REDIS.Set(ctx, cacheKey, data, dictCacheExpire)
	}

	return dict.Details, nil
}

func (s *dictionaryService) clearCache(ctx context.Context, dictType string) {
	if global.G_REDIS != nil && dictType != "" {
		global.G_REDIS.Del(ctx, dictCachePrefix+dictType)
	}
}
