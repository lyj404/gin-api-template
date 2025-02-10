package domain

import (
	"context"
	"gin-api-template/domain/entity"
)

type RefreshTokenService interface {
	GetUserByID(c context.Context, id string) (entity.User, error)
	CreateAccessToken(user *entity.User, secret string, expiry int) (accessToken string, err error)
	CreateRefreshToken(user *entity.User, secret string, expiry int) (refreshToken string, err error)
	ExtractIDFromToken(requestToken string, secret string) (string, error)
}

type LoginService interface {
	Create(c context.Context, user *entity.User) error
	GetUserByEmail(c context.Context, email string) (entity.User, error)
}
