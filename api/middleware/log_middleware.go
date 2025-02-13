package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// LoggerMiddleware 创建一个 Gin 日志中间件
func LoggerMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		// 将元信息添加到 logrus.Entry 中
		entry := logger.WithFields(logrus.Fields{
			"method":  c.Request.Method,
			"url":     c.Request.URL.Path,
			"ip":      c.ClientIP(),
			"status":  c.Writer.Status(),
			"latency": time.Since(start),
		})

		// 根据状态码设置日志级别
		logLevel := logrus.InfoLevel
		if c.Writer.Status() >= 500 {
			logLevel = logrus.ErrorLevel
		} else if c.Writer.Status() >= 400 {
			logLevel = logrus.WarnLevel
		}

		// 记录日志（具体格式由 log_formatter 处理）
		entry.Log(logLevel, "Request completed")
	}
}
