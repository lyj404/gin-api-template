//go:build wireinject
// +build wireinject

package main

import (
	"time"

	"github.com/lyj404/gin-api-template/config"
	"github.com/google/wire"
	"github.com/lyj404/gin-api-template/api/handler"
	"github.com/lyj404/gin-api-template/domain"
	"github.com/lyj404/gin-api-template/global"
	"github.com/lyj404/gin-api-template/repository"
	"github.com/lyj404/gin-api-template/service"
	"gorm.io/gorm"
)

type App struct {
	DB         *gorm.DB
	UserRepo   domain.UserRepo
	UserSvc    domain.LoginService
	TokenSvc   domain.RefreshTokenService
	UserHdlr   *handler.UserHandler
	HelloHdlr  *handler.HelloHandler
	RefreshHdlr *handler.RefreshTokenHandler
}

func InitializeApp() (*App, error) {
	wire.Build(
		providerSet,
		wire.Struct(new(App), "*"),
	)
	return &App{}, nil
}

func provideDB() *gorm.DB {
	return global.G_DB
}

func provideTimeout() time.Duration {
	return time.Duration(config.CfgTimeout.ContextTimeout) * time.Second
}

var providerSet = wire.NewSet(
	provideDB,
	provideTimeout,

	repository.NewUserRepo,
	service.NewUserService,
	service.NewRefreshTokenService,

	handler.NewUserHandler,
	handler.NewHelloHandler,
	handler.NewRefreshTokenHandler,
)
