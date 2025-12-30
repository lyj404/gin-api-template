package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lyj404/gin-api-template/config"
)

func CorsMiddleware() gin.HandlerFunc {
	// 获取默认的 CORS 配置
	corsConfig := cors.DefaultConfig()

	// 如果配置了允许的域名，使用配置的域名；否则允许所有（向后兼容）
	if len(config.CfgServer.AllowedOrigins) > 0 {
		corsConfig.AllowOrigins = config.CfgServer.AllowedOrigins
	} else {
		corsConfig.AllowAllOrigins = true
	}

	// 允许跨域请求包含凭证（如 Cookies 或 Authorization 头）
	corsConfig.AllowCredentials = true
	// 添加允许的自定义请求头
	corsConfig.AddAllowHeaders("User-Agent", "Authorization")
	return cors.New(corsConfig)
}
