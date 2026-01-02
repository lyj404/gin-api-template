//go:build wireinject
// +build wireinject

package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/lyj404/gin-api-template/api/handler"
	"github.com/lyj404/gin-api-template/api/route"
	"github.com/lyj404/gin-api-template/config"
	"github.com/lyj404/gin-api-template/domain"
	"github.com/lyj404/gin-api-template/global"
	"github.com/lyj404/gin-api-template/pkg/lib/logger"
	"github.com/lyj404/gin-api-template/repository"
	"github.com/lyj404/gin-api-template/service"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// App 应用结构体，包含所有依赖
type App struct {
	DB             *gorm.DB
	Redis          *redis.Client
	Logger         *zap.Logger
	Router         *gin.Engine
	UserRepo       domain.UserRepo
	UserSvc        domain.LoginService
	TokenSvc       domain.RefreshTokenService
	UserHdlr       *handler.UserHandler
	HelloHdlr      *handler.HelloHandler
	RefreshHdlr    *handler.RefreshTokenHandler
	RegisterRoutes func()
}

// InitializeApp 初始化应用，使用 Wire 自动生成依赖注入代码
func InitializeApp() (*App, error) {
	wire.Build(
		providerSet,
		wire.Struct(new(App), "*"),
	)
	return &App{}, nil
}

// provideDB 提供 DB 实例
func provideDB() *gorm.DB {
	return global.G_DB
}

// provideRedis 提供 Redis 实例
func provideRedis() *redis.Client {
	return global.G_REDIS
}

// provideLogger 提供 Logger 实例
func provideLogger() *zap.Logger {
	return logger.InitZapLogger()
}

// provideRouter 提供 Router 实例，初始化路由和中间件
func provideRouter(timeout time.Duration, logger *zap.Logger) *gin.Engine {
	return route.SetUp(timeout, logger)
}

// provideRouteRegistration 提供路由注册函数
func provideRouteRegistration(
	router *gin.Engine,
	userHdlr *handler.UserHandler,
	helloHdlr *handler.HelloHandler,
	refreshTokenHdlr *handler.RefreshTokenHandler,
) func() {
	return func() {
		// 注册公共路由
		publicGroup := router.Group("")
		route.NewUserRouter(userHdlr, refreshTokenHdlr, publicGroup)

		// 注册受保护的路由
		protectedGroup := router.Group("")
		protectedGroup.Use(route.JwtAuthMiddleware())
		route.NewTestRouter(helloHdlr, protectedGroup)
	}
}

// provideTimeout 提供超时时间
func provideTimeout() time.Duration {
	return time.Duration(config.CfgTimeout.ContextTimeout) * time.Second
}

// providerSet 依赖注入提供者集合
var providerSet = wire.NewSet(
	provideDB,
	provideRedis,
	provideLogger,
	provideRouter,
	provideRouteRegistration,
	provideTimeout,

	// Repository 层
	repository.NewUserRepo,

	// Service 层
	service.NewUserService,
	service.NewRefreshTokenService,

	// Handler 层
	handler.NewUserHandler,
	handler.NewHelloHandler,
	handler.NewRefreshTokenHandler,
)
