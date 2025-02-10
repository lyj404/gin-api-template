package route

import (
	"gin-api-template/api/handler"
	"gin-api-template/bootstrap"
	"gin-api-template/repo"
	"gin-api-template/service"
	"time"

	"github.com/gin-gonic/gin"
)

func NewUserRouter(timeout time.Duration, app bootstrap.Application, group *gin.RouterGroup) {
	userRepo := repo.NewUserRepo(app.Db)
	userHandler := &handler.UserHandler{
		UserService:         service.NewUserService(userRepo, timeout),
		RefreshTokenUseCase: service.NewRefreshTokenService(userRepo, timeout),
	}
	refreshTokenHandler := &handler.RefreshTokenHandler{
		RefreshTokenService: service.NewRefreshTokenService(userRepo, timeout),
	}
	group.POST("/login", userHandler.Login)
	group.POST("/signup", userHandler.Signup)
	group.POST("/refresh", refreshTokenHandler.RefreshToken)
}
