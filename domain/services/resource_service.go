package services

import "github.com/lyj404/gin-api-template/domain/entity"

type ResourceService interface {
	CreateResource(resource *entity.Resource, operatorID uint64) error
	UpdateResource(resource *entity.Resource, operatorID uint64) error
	DeleteResource(id uint64, operatorID uint64) error
	GetResourceByID(id uint64) (*entity.Resource, error)
	GetAllResources() ([]entity.Resource, error)
}
