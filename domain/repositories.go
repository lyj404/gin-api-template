package domain

import (
	"context"

	"github.com/lyj404/gin-api-template/domain/entity"
)

type UserRepo interface {
	Create(c context.Context, user *entity.User) error
	GetByEmail(c context.Context, email string) (entity.User, error)
	GetByID(c context.Context, id string) (entity.User, error)
}
