package result

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseResult[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    *T     `json:"data,omitempty"`
}

// ErrorResponse 返回错误响应
func ErrorResponse(ctx *gin.Context, code int, message string) {
	ctx.JSON(code, ResponseResult[string]{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// SuccessResponse 返回成功响应
func SuccessResponse[T any](ctx *gin.Context, message string, data *T) {
	ctx.JSON(http.StatusOK, ResponseResult[T]{
		Code:    http.StatusOK,
		Message: message,
		Data:    data,
	})
}
