// Package middleware 提供中间件功能
package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORS 跨域中间件
// 允许前端从不同域名访问API
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 允许所有来源（生产环境建议设置具体域名）
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		// 允许的请求方法
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// 允许的请求头
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		// 允许携带认证信息
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// 预检请求的缓存时间（秒）
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}
