package repo

import (
	"context"
	"gin-api-template/domain"

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

func (u *userRepo) Create(c context.Context, user *domain.User) error {
	return u.database.WithContext(c).Create(user).Error
}

func (u *userRepo) GetByEmail(c context.Context, email string) (domain.User, error) {
	var user domain.User
	err := u.database.WithContext(c).Where("email = ?", email).First(&user).Error
	return user, err
}

func (u *userRepo) GetByID(c context.Context, id string) (domain.User, error) {
	var user domain.User
	if err := u.database.WithContext(c).First(&user, id).Error; err != nil {
		return domain.User{}, err
	}
	return user, nil
}
