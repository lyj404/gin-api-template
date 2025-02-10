package middleware

import (
	"gin-api-template/domain/result"
	"gin-api-template/internal/tokenutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取 Authorization 字段
		authHeader := c.Request.Header.Get("Authorization")
		// 按空格分割 Authorization 字段，通常格式为 "Bearer <token>"
		t := strings.Split(authHeader, " ")
		// 检查分割后的结果长度是否为 2，确保格式正确「方案选单」
		if len(t) == 2 {
			// 提取令牌部分
			authToken := t[1]
			// 验证令牌是否被授权
			authorized, err := tokenutil.IsAuthorized(authToken, secret)
			if authorized {
				// 如果授权成功，从令牌中提取用户 ID
				userID, err := tokenutil.ExtractIDFromToken(authToken, secret)
				if err != nil {
					// 如果提取用户 ID 时出错，返回未授权错误
					c.JSON(http.StatusUnauthorized, result.ErrorResponse{Message: err.Error()})
					c.Abort()
					return
				}
				// 将用户 ID 设置到上下文中，供后续处理使用
				c.Set("x-user-id", userID)
				// 继续处理请求
				c.Next()
				return
			}
			// 如果授权失败，返回未授权错误
			c.JSON(http.StatusUnauthorized, result.ErrorResponse{Message: err.Error()})
			c.Abort()
			return
		}
		// 如果 Authorization 格式不正确，返回未授权错误
		c.JSON(http.StatusUnauthorized, result.ErrorResponse{Message: "Not authorized"})
		c.Abort()
	}
}
