package route

import (
	"gin-api-template/api/handler"
	"time"

	"github.com/gin-gonic/gin"
)

func NewTestRouter(timeout time.Duration, group *gin.RouterGroup) {
	testHandler := &handler.HelloHandler{}
	group.GET("/hello", testHandler.HelloRequest)
}
