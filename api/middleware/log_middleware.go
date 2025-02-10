package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// LoggerMiddleware 创建一个 Gin 日志中间件
func LoggerMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		start := time.Now()

		// 记录请求的基本信息
		logger.WithFields(logrus.Fields{
			"method": c.Request.Method,   // 请求方法
			"url":    c.Request.URL.Path, // 请求路径
			"ip":     c.ClientIP(),       // 客户端 IP
		}).Info("Request started")

		// 处理请求
		c.Next()

		// 根据状态码设置日志级别
		logLevel := logrus.InfoLevel
		if c.Writer.Status() >= 500 {
			logLevel = logrus.ErrorLevel
		} else if c.Writer.Status() >= 400 {
			logLevel = logrus.WarnLevel
		}

		// 记录请求完成的信息
		logger.WithFields(logrus.Fields{
			"status":  c.Writer.Status(), // 响应状态码
			"latency": time.Since(start), // 请求处理耗时
		}).Log(logLevel, "Request completed")
	}
}
