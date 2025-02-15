package route

import (
	"gin-api-template/api/middleware"
	"gin-api-template/bootstrap"
	"gin-api-template/config"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func SetUp(timeout time.Duration, app bootstrap.Application, logger *logrus.Logger) *gin.Engine {
	router := gin.New()
	// 使用自定义日志
	router.Use(middleware.LoggerMiddleware(logger))

	publicRouter := router.Group("")
	// 不需要鉴权的路由
	NewUserRouter(timeout, app, publicRouter)

	protectedRouter := router.Group("")
	// 需要鉴权的路由
	protectedRouter.Use(middleware.JwtAuthMiddleware(config.CfgToken.AccessTokenSecret))
	NewTestRouter(timeout, app, protectedRouter)

	return router
}
