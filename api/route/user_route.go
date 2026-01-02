package route

import (
	"github.com/lyj404/gin-api-template/api/handler"

	"github.com/gin-gonic/gin"
)

func NewUserRouter(userHdlr *handler.UserHandler, refreshTokenHdlr *handler.RefreshTokenHandler, group *gin.RouterGroup) {
	group.POST("/login", userHdlr.Login)
	group.POST("/signup", userHdlr.Signup)
	group.POST("/refresh-token", refreshTokenHdlr.RefreshToken)
	group.GET("/captcha", userHdlr.GenerateMathCaptcha)
}
