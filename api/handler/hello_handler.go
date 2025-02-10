package handler

import (
	"gin-api-template/domain/result"

	"github.com/gin-gonic/gin"
)

type HelloHandler struct {
}

func (handler *HelloHandler) TestRequest(c *gin.Context) {
	result.SuccessResponse[string](c, "test", nil)
}

func (handler *HelloHandler) HelloRequest(c *gin.Context) {
	result.SuccessResponse[string](c, "Hello World", nil)
}
