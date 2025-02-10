package service

import (
	"context"
	"gin-api-template/domain"
	"time"
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

func (u *userService) GetUserByEmail(c context.Context, email string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeOut)
	defer cancel()
	return u.repo.GetByEmail(ctx, email)
}

func (u *userService) Create(c context.Context, user *domain.User) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeOut)
	defer cancel()
	return u.repo.Create(ctx, user)
}
