package handler

import (
	"gin-api-template/internal/redisutil"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type TestHandler struct {
	Cache *redisutil.RedisClient
}

func (tc *TestHandler) TestRequest(c *gin.Context) {
	tc.Cache.Set(c, "test", "Hello World", 1*time.Minute)
	var val string
	tc.Cache.Get(c, "test", &val)
	log.Println("val:", val, "==========")
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello World",
	})
}
