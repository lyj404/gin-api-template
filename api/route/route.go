package route

import (
	"os"
	"time"

	"github.com/lyj404/gin-api-template/api/middleware"
	"github.com/lyj404/gin-api-template/config"
	_ "github.com/lyj404/gin-api-template/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

func SetUp(timeout time.Duration, logger *zap.Logger) *gin.Engine {
	// 设置gin运行模式
	gin.SetMode(config.CfgServer.Mode)
	router := gin.New()

	// 使用TraceID中间件（必须在最前面）
	router.Use(middleware.TraceIDMiddleware())

	// 使用自定义Recovery中间件
	router.Use(middleware.RecoveryMiddleware(logger))

	// 使用自定义日志
	router.Use(middleware.LoggerMiddleware(logger))

	// 使用错误处理中间件
	router.Use(middleware.ErrorHandlerMiddleware(logger))

	// 使用CORS中间件
	router.Use(middleware.CorsMiddleware())

	// 使用限流中间件（如果配置了限流）
	if config.CfgServer.RateLimit > 0 {
		if config.CfgRedis.Enabled {
			router.Use(middleware.RateLimitMiddlewareWithRedis(config.CfgServer.RateLimit))
		} else {
			router.Use(middleware.RateLimitMiddleware(config.CfgServer.RateLimit))
		}
	}

	// 设置swagger路由
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}

// JwtAuthMiddleware JWT 鉴权中间件
func JwtAuthMiddleware() gin.HandlerFunc {
	return middleware.JwtAuthMiddleware(config.CfgToken.AccessTokenSecret)
}

// SetupFrontend 设置前端静态文件服务与 SPA 回退
func SetupFrontend(router *gin.Engine) {
	// 确定前端 dist 目录路径
	wd, _ := os.Getwd()
	distPath := wd + "/web/dist"

	// 检查 dist 目录是否存在
	if _, err := os.Stat(distPath); os.IsNotExist(err) {
		// 尝试向上一级查找
		distPath = wd + "/../web/dist"
		if _, err := os.Stat(distPath); os.IsNotExist(err) {
			return // 前端未构建，跳过
		}
	}

	// 静态文件路由 - 优先匹配具体文件
	router.Use(func(c *gin.Context) {
		// 如果是 API 路径或 swagger，跳过静态文件处理
		if c.Request.URL.Path == "/swagger/index.html" || c.Request.URL.Path == "/swagger/ui/index.html" {
			c.Next()
			return
		}
		// 检查请求路径是否为前端路由（非 API 路径）
		if !isAPIPath(c.Request.URL.Path) {
			// 如果是前端路由（没有文件扩展名），返回 index.html
			if !hasFileExtension(c.Request.URL.Path) {
				c.Header("Content-Type", "text/html; charset=utf-8")
				c.File(distPath + "/index.html")
				c.Abort()
				return
			}
		}
		c.Next()
	})

	router.Use(func(c *gin.Context) {
		if isAPIPath(c.Request.URL.Path) || c.Request.URL.Path == "/swagger" {
			c.Next()
			return
		}

		// 尝试提供静态文件
		filePath := distPath + c.Request.URL.Path
		if _, err := os.Stat(filePath); err == nil {
			// 根据文件扩展名设置正确的 Content-Type
			contentType := getContentType(c.Request.URL.Path)
			c.Header("Content-Type", contentType)
			c.File(filePath)
			c.Abort()
			return
		}

		// 如果文件不存在，返回 index.html（SPA 回退）
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.File(distPath + "/index.html")
		c.Abort()
	})
}

// isAPIPath 检查路径是否为 API 路径
func isAPIPath(path string) bool {
	return len(path) >= 5 && path[:5] == "/api/"
}

// hasFileExtension 检查路径是否有文件扩展名
func hasFileExtension(path string) bool {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '.' {
			return true
		}
		if path[i] == '/' {
			return false
		}
	}
	return false
}

// getContentType 根据文件扩展名返回正确的 Content-Type
func getContentType(path string) string {
	ext := ""
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '.' {
			ext = path[i:]
			break
		}
		if path[i] == '/' {
			break
		}
	}

	switch ext {
	case ".html", ".htm":
		return "text/html; charset=utf-8"
	case ".css":
		return "text/css; charset=utf-8"
	case ".js":
		return "application/javascript; charset=utf-8"
	case ".json":
		return "application/json; charset=utf-8"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".svg":
		return "image/svg+xml"
	case ".ico":
		return "image/x-icon"
	default:
		return "text/plain; charset=utf-8"
	}
}
