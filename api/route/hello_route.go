package route

import (
	"github.com/lyj404/gin-api-template/api/handler"

	"github.com/gin-gonic/gin"
)

func NewTestRouter(helloHdlr *handler.HelloHandler, group *gin.RouterGroup) {
	group.GET("/hello", helloHdlr.HelloRequest)
}
