package middleware

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lyj404/gin-api-template/domain/result"
	"github.com/lyj404/gin-api-template/global"
)

// RateLimiter 限流器结构
type RateLimiter struct {
	mu      sync.Mutex
	clients map[string]*ClientInfo
}

type ClientInfo struct {
	Tokens     int
	LastUpdate time.Time
}

var limiter *RateLimiter

// InitRateLimiter 初始化限流器
func InitRateLimiter() {
	limiter = &RateLimiter{
		clients: make(map[string]*ClientInfo),
	}
}

// RateLimitMiddleware 基于内存的限流中间件
func RateLimitMiddleware(requestsPerMinute int) gin.HandlerFunc {
	if limiter == nil {
		InitRateLimiter()
	}

	return func(c *gin.Context) {
		ip := c.ClientIP()

		limiter.mu.Lock()
		client, exists := limiter.clients[ip]
		if !exists {
			client = &ClientInfo{
				Tokens:     requestsPerMinute - 1,
				LastUpdate: time.Now(),
			}
			limiter.clients[ip] = client
			limiter.mu.Unlock()
			c.Next()
			return
		}

		now := time.Now()
		elapsed := now.Sub(client.LastUpdate)
		client.LastUpdate = now

		tokensToAdd := int(elapsed.Minutes() * float64(requestsPerMinute))
		client.Tokens += tokensToAdd
		if client.Tokens > requestsPerMinute {
			client.Tokens = requestsPerMinute
		}

		if client.Tokens <= 0 {
			limiter.mu.Unlock()
			result.ErrorResponse(c, http.StatusTooManyRequests, "Too many requests, please try again later")
			c.Abort()
			return
		}

		client.Tokens--
		limiter.mu.Unlock()
		c.Next()
	}
}

// RateLimitMiddlewareWithRedis 基于 Redis 的限流中间件
func RateLimitMiddlewareWithRedis(requestsPerMinute int) gin.HandlerFunc {
	if global.G_REDIS == nil {
		return RateLimitMiddleware(requestsPerMinute)
	}

	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := "rate_limit:" + ip

		pipe := global.G_REDIS.Pipeline()
		current := pipe.Get(context.Background(), key)
		pipe.Expire(context.Background(), key, time.Minute)

		_, err := pipe.Exec(context.Background())
		if err != nil {
			c.Next()
			return
		}

		count, _ := current.Int()
		if count >= requestsPerMinute {
			result.ErrorResponse(c, http.StatusTooManyRequests, "Too many requests, please try again later")
			c.Abort()
			return
		}

		global.G_REDIS.Incr(context.Background(), key)
		c.Next()
	}
}
