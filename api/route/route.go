package route

import (
	"gin-api-template/api/middleware"
	"gin-api-template/bootstrap"
	"gin-api-template/config"
	"time"

	"github.com/gin-gonic/gin"
)

func SetUp(timeout time.Duration, app bootstrap.Application, gin *gin.Engine) {
	publicRouter := gin.Group("")
	// 不需要鉴权的路由
	NewUserRouter(timeout, app, publicRouter)

	protectedRouter := gin.Group("")
	// 需要鉴权的路由
	protectedRouter.Use(middleware.JwtAuthMiddleware(config.AccessTokenSecret))
	NewTestRouter(timeout, app, protectedRouter)
}
