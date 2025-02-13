package result

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ResponseResult 是通用响应结构体
// @Description 统一的API响应格式
type ResponseResult[T any] struct {
	// @Description 响应状态码
	Code int `json:"code"`
	// @Description 响应信息
	Message string `json:"message"`
	// @Description 响应数据
	Data *T `json:"data,omitempty"`
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

// SimpleSuccessResponse 返回简单的成功响应
func SimpleSuccessResponse(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusOK, ResponseResult[any]{
		Code:    http.StatusOK,
		Message: message,
	})
}
