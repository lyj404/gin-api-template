package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoggerMiddleware 创建一个 Gin 日志中间件
func LoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		// 计算请求耗时
		latency := time.Since(start)

		// 创建日志字段
		fields := []zap.Field{
			zap.String("method", c.Request.Method),
			zap.String("url", c.Request.URL.Path),
			zap.String("ip", c.ClientIP()),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", latency),
		}

		// 根据状态码设置日志级别
		switch {
		case c.Writer.Status() >= 500:
			logger.Error("Request completed", fields...)
		case c.Writer.Status() >= 400:
			logger.Warn("Request completed", fields...)
		default:
			logger.Info("Request completed", fields...)
		}
	}
}
