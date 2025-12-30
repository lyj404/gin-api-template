package route

import (
	"time"

	"github.com/lyj404/gin-api-template/api/middleware"
	"github.com/lyj404/gin-api-template/config"
	_ "github.com/lyj404/gin-api-template/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

func SetUp(timeout time.Duration, logger *zap.Logger) *gin.Engine {
	// 设置gin运行模式
	gin.SetMode(config.CfgServer.Mode)
	router := gin.New()

	// 使用TraceID中间件（必须在最前面）
	router.Use(middleware.TraceIDMiddleware())

	// 使用自定义Recovery中间件
	router.Use(middleware.RecoveryMiddleware(logger))

	// 使用自定义日志
	router.Use(middleware.LoggerMiddleware(logger))

	// 使用错误处理中间件
	router.Use(middleware.ErrorHandlerMiddleware(logger))

	// 使用CORS中间件
	router.Use(middleware.CorsMiddleware())

	// 使用限流中间件（如果配置了限流）
	if config.CfgServer.RateLimit > 0 {
		if config.CfgRedis.Enabled {
			router.Use(middleware.RateLimitMiddlewareWithRedis(config.CfgServer.RateLimit))
		} else {
			router.Use(middleware.RateLimitMiddleware(config.CfgServer.RateLimit))
		}
	}

	// 设置swagger路由
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	publicRouter := router.Group("")
	// 不需要鉴权的路由
	NewUserRouter(timeout, publicRouter)

	protectedRouter := router.Group("")
	// 需要鉴权的路由
	protectedRouter.Use(middleware.JwtAuthMiddleware(config.CfgToken.AccessTokenSecret))
	NewTestRouter(timeout, protectedRouter)

	return router
}
