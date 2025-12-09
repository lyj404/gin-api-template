package route

import (
	"time"

	"github.com/lyj404/gin-api-template/api/handler"

	"github.com/gin-gonic/gin"
)

func NewTestRouter(timeout time.Duration, group *gin.RouterGroup) {
	testHandler := &handler.HelloHandler{}
	group.GET("/hello", testHandler.HelloRequest)
}
