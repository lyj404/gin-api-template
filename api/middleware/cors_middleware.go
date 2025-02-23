package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CorsMiddleware() gin.HandlerFunc {
	// 获取默认的 CORS 配置
	config := cors.DefaultConfig()
	// 允许所有来源（*）的跨域请求
	config.AllowAllOrigins = true
	// 允许跨域请求包含凭证（如 Cookies 或 Authorization 头）
	config.AllowCredentials = true
	// 添加允许的自定义请求头
	config.AddAllowHeaders("User-Agent", "Authorization")
	return cors.New(config)
}
