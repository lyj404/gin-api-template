package route

import (
	"gin-api-template/api/handler"
	"gin-api-template/bootstrap"
	"gin-api-template/internal/redisutil"
	"time"

	"github.com/gin-gonic/gin"
)

func NewTestRouter(timeout time.Duration, app bootstrap.Application, group *gin.RouterGroup) {
	cache := redisutil.NewRedisClient(app.Cache)
	testHandler := &handler.TestHandler{
		Cache: cache,
	}
	group.GET("/test", testHandler.TestRequest)
}
