package repositories

import "github.com/lyj404/gin-api-template/domain/entity"

type ResourceRepository interface {
	Create(resource *entity.Resource) error
	Update(resource *entity.Resource) error
	Delete(id uint) error
	GetByID(id uint) (*entity.Resource, error)
	GetAll() ([]entity.Resource, error)
}
