package handler

import (
	"gin-api-template/domain/result"

	"github.com/gin-gonic/gin"
)

type HelloHandler struct {
}

func (handler *HelloHandler) TestRequest(c *gin.Context) {
	result.SimpleSuccessResponse(c, "test")
}

func (handler *HelloHandler) HelloRequest(c *gin.Context) {
	result.SimpleSuccessResponse(c, "Hello World")
}
