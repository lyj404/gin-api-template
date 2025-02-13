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
		if authHeader == "" {
			result.ErrorResponse(c, http.StatusUnauthorized, "token not found")
			c.Abort()
			return
		}

		// 分割 Authorization 字段
		var authToken string
		if strings.HasPrefix(authHeader, "Bearer ") {
			// 如果包含 Bearer 前缀，则提取令牌部分
			authToken = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			// 如果不包含 Bearer 前缀，直接使用整个字符串作为令牌
			authToken = authHeader
		}

		// 验证令牌是否有效及过期
		_, err := tokenutil.ExtractIDFromToken(authToken, secret)
		if err != nil {
			// 根据错误类型返回相应的错误信息
			if err.Error() == "token is expired" {
				result.ErrorResponse(c, http.StatusUnauthorized, "Token is expired")
			} else {
				result.ErrorResponse(c, http.StatusUnauthorized, err.Error())
			}
			c.Abort()
			return
		}

		// 继续处理请求
		c.Next()
	}
}
