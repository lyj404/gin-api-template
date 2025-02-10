package route

import (
	"gin-api-template/api/handler"
	"gin-api-template/bootstrap"
	"time"

	"github.com/gin-gonic/gin"
)

func NewTestRouter(timeout time.Duration, app bootstrap.Application, group *gin.RouterGroup) {
	testHandler := &handler.HelloHandler{}
	group.GET("/test", testHandler.TestRequest)
	group.GET("/hello", testHandler.HelloRequest)
}
