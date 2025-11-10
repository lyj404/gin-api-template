package route

import (
	"gin-api-template/api/middleware"
	"gin-api-template/config"
	_ "gin-api-template/docs"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

func SetUp(timeout time.Duration, logger *zap.Logger) *gin.Engine {
	// 设置gin运行模式
	gin.SetMode(config.CfgServer.Mode)
	router := gin.New()

	// 使用Recovery中间件
	router.Use(gin.Recovery())

	// 使用自定义日志
	router.Use(middleware.LoggerMiddleware(logger))

	// 使用CORS中间件
	router.Use(middleware.CorsMiddleware())

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
