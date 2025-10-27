// Package utils 提供工具函数
package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims JWT声明结构
type JWTClaims struct {
	UserID   int64  `json:"user_id"`  // 用户ID
	Username string `json:"username"` // 用户名
	jwt.RegisteredClaims
}

// 密钥（生产环境应该从配置文件读取）
var jwtSecret = []byte("your-secret-key-change-in-production")

// GenerateToken 生成JWT token
// 参数：userID 用户ID, username 用户名
// 返回：token字符串和错误
func GenerateToken(userID int64, username string) (string, error) {
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24小时过期
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken 解析JWT token
// 参数：tokenString token字符串
// 返回：claims和错误
func ParseToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
