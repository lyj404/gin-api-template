package route

import (
	"time"

	"github.com/lyj404/gin-api-template/api/handler"
	"github.com/lyj404/gin-api-template/global"
	"github.com/lyj404/gin-api-template/repository"
	"github.com/lyj404/gin-api-template/service"

	"github.com/gin-gonic/gin"
)

func NewUserRouter(timeout time.Duration, group *gin.RouterGroup) {
	userRepo := repository.NewUserRepo(global.G_DB)
	userHandler := &handler.UserHandler{
		UserService:         service.NewUserService(userRepo, timeout),
		RefreshTokenUseCase: service.NewRefreshTokenService(userRepo, timeout),
	}
	refreshTokenHandler := &handler.RefreshTokenHandler{
		RefreshTokenService: service.NewRefreshTokenService(userRepo, timeout),
	}
	group.POST("/login", userHandler.Login)
	group.POST("/signup", userHandler.Signup)
	group.POST("/refresh-token", refreshTokenHandler.RefreshToken)
	group.GET("/captcha", userHandler.GenerateMathCaptcha)
}
