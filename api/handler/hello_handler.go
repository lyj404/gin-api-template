package handler

import (
	"github.com/lyj404/gin-api-template/domain/result"

	"github.com/gin-gonic/gin"
)

type HelloHandler struct {
}

// @Summary 测试
// @Description 返回一个简单的消息
// @Tags hello
// @Accept json
// @Produce json
// @Success 200 {object} result.ResponseResult[string] "成功响应"
// @Failure 500 {object} result.ResponseResult[string] "服务器错误"
// @Router /hello [get]
func (handler *HelloHandler) HelloRequest(c *gin.Context) {
	result.SimpleSuccessResponse(c, "Hello World")
}
