package repository

import (
	"github.com/lyj404/gin-api-template/domain/entity"
	"github.com/lyj404/gin-api-template/domain/repositories"
	"github.com/lyj404/gin-api-template/global"
)

type resourceRepository struct{}

func NewResourceRepository() repositories.ResourceRepository {
	return &resourceRepository{}
}

func (r *resourceRepository) Create(resource *entity.Resource) error {
	return global.G_DB.Create(resource).Error
}

func (r *resourceRepository) Update(resource *entity.Resource) error {
	return global.G_DB.Save(resource).Error
}

func (r *resourceRepository) Delete(id uint) error {
	return global.G_DB.Delete(&entity.Resource{}, id).Error
}

func (r *resourceRepository) GetByID(id uint) (*entity.Resource, error) {
	var resource entity.Resource
	err := global.G_DB.First(&resource, id).Error
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

func (r *resourceRepository) GetAll() ([]entity.Resource, error) {
	var resources []entity.Resource
	err := global.G_DB.Find(&resources).Error
	return resources, err
}
