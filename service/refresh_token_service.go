package service

import (
	"context"
	"time"

	"github.com/lyj404/gin-api-template/domain"
	"github.com/lyj404/gin-api-template/domain/entity"
	"github.com/lyj404/gin-api-template/internal/tokenutil"
)

type refreshTokenService struct {
	userRepo       domain.UserRepo
	contextTimeout time.Duration
}

func NewRefreshTokenService(userRepo domain.UserRepo, timeout time.Duration) domain.RefreshTokenService {
	return &refreshTokenService{
		userRepo:       userRepo,
		contextTimeout: timeout,
	}
}

func (rtu *refreshTokenService) GetUserByID(c context.Context, email string) (entity.User, error) {
	ctx, cancel := context.WithTimeout(c, rtu.contextTimeout)
	defer cancel()
	return rtu.userRepo.GetByID(ctx, email)
}

func (rtu *refreshTokenService) CreateAccessToken(user *entity.User, secret string, expiry int) (accessToken string, err error) {
	return tokenutil.CreateAccessToken(user, secret, expiry)
}

func (rtu *refreshTokenService) CreateRefreshToken(user *entity.User, secret string, expiry int) (refreshToken string, err error) {
	return tokenutil.CreateRefreshToken(user, secret, expiry)
}

func (rtu *refreshTokenService) ExtractIDFromToken(requestToken string, secret string) (string, error) {
	return tokenutil.ExtractIDFromToken(requestToken, secret)
}
