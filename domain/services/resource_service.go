package services

import "github.com/lyj404/gin-api-template/domain/entity"

type ResourceService interface {
	CreateResource(resource *entity.Resource, operatorID uint) error
	UpdateResource(resource *entity.Resource, operatorID uint) error
	DeleteResource(id uint, operatorID uint) error
	GetResourceByID(id uint) (*entity.Resource, error)
	GetAllResources() ([]entity.Resource, error)
}
