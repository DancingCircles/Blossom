// Package middleware 中间件包
package middleware

import (
	"net/http"
	"sync"
	"web_app/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter 全局限流器
var (
	globalLimiter = rate.NewLimiter(100, 200) // 每秒100个请求，桶容量200
	userLimiters  = make(map[int64]*rate.Limiter)
	mu            sync.Mutex
)

// RateLimit 全局限流中间件
// 使用令牌桶算法，限制每秒请求数
func RateLimit(rps int, burst int) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Limit(rps), burst)

	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, models.NewErrorResponse(models.CodeTooManyRequests, "请求过于频繁，请稍后再试"))
			c.Abort()
			return
		}
		c.Next()
	}
}

// UserRateLimit 基于用户ID的限流中间件
// 针对已登录用户进行限流，每个用户独立计算
func UserRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户ID（由JWT中间件注入）
		userIDVal, exists := c.Get("user_id")
		if !exists {
			// 未登录用户使用全局限流器
			if !globalLimiter.Allow() {
				c.JSON(http.StatusTooManyRequests, models.NewErrorResponse(models.CodeTooManyRequests, "请求过于频繁，请稍后再试"))
				c.Abort()
				return
			}
			c.Next()
			return
		}

		userID, ok := userIDVal.(int64)
		if !ok {
			c.JSON(http.StatusInternalServerError, models.NewErrorResponse(models.CodeServerError, "用户ID格式错误"))
			c.Abort()
			return
		}

		// 获取或创建该用户的限流器
		mu.Lock()
		limiter, exists := userLimiters[userID]
		if !exists {
			// 为该用户创建新的限流器：每秒50个请求，桶容量100
			limiter = rate.NewLimiter(50, 100)
			userLimiters[userID] = limiter
		}
		mu.Unlock()

		// 检查是否超出限流
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, models.NewErrorResponse(models.CodeTooManyRequests, "您的请求过于频繁，请稍后再试"))
			c.Abort()
			return
		}

		c.Next()
	}
}

// IPRateLimit 基于IP的限流中间件
// 针对IP地址进行限流，防止单个IP恶意请求
func IPRateLimit(rps int, burst int) gin.HandlerFunc {
	ipLimiters := make(map[string]*rate.Limiter)
	mu := sync.Mutex{}

	return func(c *gin.Context) {
		ip := c.ClientIP()

		// 获取或创建该IP的限流器
		mu.Lock()
		limiter, exists := ipLimiters[ip]
		if !exists {
			limiter = rate.NewLimiter(rate.Limit(rps), burst)
			ipLimiters[ip] = limiter
		}
		mu.Unlock()

		// 检查是否超出限流
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, models.NewErrorResponse(models.CodeTooManyRequests, "请求过于频繁，请稍后再试"))
			c.Abort()
			return
		}

		c.Next()
	}
}
