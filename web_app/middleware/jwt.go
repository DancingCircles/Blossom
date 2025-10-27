// Package middleware 提供中间件功能
package middleware

import (
	"net/http"
	"strings"
	"web_app/models"
	"web_app/utils"

	"github.com/gin-gonic/gin"
)

// GenerateToken 生成JWT token（包装utils中的函数）
func GenerateToken(userID int64, username string) (string, error) {
	return utils.GenerateToken(userID, username)
}

// JWTAuth JWT认证中间件
// 用于验证请求头中的JWT token
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从请求头获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.CodeUnauthorized, "请先登录"))
			c.Abort()
			return
		}

		// 2. 检查token格式：Bearer <token>
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.CodeUnauthorized, "token格式错误"))
			c.Abort()
			return
		}

		// 3. 解析和验证JWT token
		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.CodeUnauthorized, "token无效或已过期"))
			c.Abort()
			return
		}

		// 4. 将用户信息存入context，供后续处理使用
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		// 继续处理请求
		c.Next()
	}
}

// OptionalJWTAuth 可选的JWT认证中间件
// 如果有token就验证，没有token也可以继续
func OptionalJWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// 没有token，继续处理
			c.Next()
			return
		}

		// 有token，尝试验证
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			// 解析token并设置用户信息到context
			if claims, err := utils.ParseToken(parts[1]); err == nil {
				c.Set("user_id", claims.UserID)
				c.Set("username", claims.Username)
			}
		}

		c.Next()
	}
}
