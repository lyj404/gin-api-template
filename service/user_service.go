package service

import (
	"context"
	"time"

	"github.com/lyj404/gin-api-template/domain"
	"github.com/lyj404/gin-api-template/domain/entity"
)

type userService struct {
	repo           domain.UserRepo
	contextTimeOut time.Duration
}

func NewUserService(repository domain.UserRepo, timeout time.Duration) domain.LoginService {
	return &userService{
		repo:           repository,
		contextTimeOut: timeout,
	}
}

func (u *userService) GetUserByEmail(c context.Context, email string) (entity.User, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeOut)
	defer cancel()
	return u.repo.GetByEmail(ctx, email)
}

func (u *userService) Create(c context.Context, user *entity.User) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeOut)
	defer cancel()
	return u.repo.Create(ctx, user)
}
