package repository

import (
	"context"
	"gin-api-template/domain"
	"gin-api-template/domain/entity"

	"gorm.io/gorm"
)

type userRepo struct {
	database *gorm.DB
}

func NewUserRepo(db *gorm.DB) domain.UserRepo {
	return &userRepo{
		database: db,
	}
}

func (u *userRepo) Create(c context.Context, user *entity.User) error {
	return u.database.WithContext(c).Create(user).Error
}

func (u *userRepo) GetByEmail(c context.Context, email string) (entity.User, error) {
	var user entity.User
	err := u.database.WithContext(c).Where("email = ?", email).First(&user).Error
	return user, err
}

func (u *userRepo) GetByID(c context.Context, id string) (entity.User, error) {
	var user entity.User
	if err := u.database.WithContext(c).First(&user, id).Error; err != nil {
		return entity.User{}, err
	}
	return user, nil
}
