// Package utils_test 提供工具函数测试
package utils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

// TestGenerateToken 测试JWT token生成
func TestGenerateToken(t *testing.T) {
	userID := int64(12345)
	username := "testuser"

	// 测试token生成
	token, err := GenerateToken(userID, username)
	assert.NoError(t, err, "token生成应该成功")
	assert.NotEmpty(t, token, "生成的token不应该为空")

	// 验证token格式（JWT token应该包含两个点）
	count := 0
	for _, char := range token {
		if char == '.' {
			count++
		}
	}
	assert.Equal(t, 2, count, "JWT token应该包含2个点")
}

// TestParseToken 测试JWT token解析
func TestParseToken(t *testing.T) {
	userID := int64(12345)
	username := "testuser"

	// 生成token
	token, err := GenerateToken(userID, username)
	assert.NoError(t, err)

	// 解析token
	claims, err := ParseToken(token)
	assert.NoError(t, err, "token解析应该成功")
	assert.NotNil(t, claims, "claims不应该为空")
	assert.Equal(t, userID, claims.UserID, "用户ID应该匹配")
	assert.Equal(t, username, claims.Username, "用户名应该匹配")
}

// TestParseInvalidToken 测试解析无效token
func TestParseInvalidToken(t *testing.T) {
	invalidToken := "invalid.token.string"

	// 解析无效token
	claims, err := ParseToken(invalidToken)
	assert.Error(t, err, "解析无效token应该失败")
	assert.Nil(t, claims, "claims应该为空")
}

// TestParseEmptyToken 测试解析空token
func TestParseEmptyToken(t *testing.T) {
	emptyToken := ""

	// 解析空token
	claims, err := ParseToken(emptyToken)
	assert.Error(t, err, "解析空token应该失败")
	assert.Nil(t, claims, "claims应该为空")
}

// TestParseExpiredToken 测试解析过期token（需要手动创建过期token）
func TestParseExpiredToken(t *testing.T) {
	// 创建一个已过期的token
	claims := JWTClaims{
		UserID:   12345,
		Username: "testuser",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // 1小时前过期
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	assert.NoError(t, err)

	// 解析过期token
	parsedClaims, err := ParseToken(tokenString)
	assert.Error(t, err, "解析过期token应该失败")
	assert.Nil(t, parsedClaims, "claims应该为空")
}

// TestTokenRoundTrip 测试token生成和解析的完整流程
func TestTokenRoundTrip(t *testing.T) {
	testCases := []struct {
		name     string
		userID   int64
		username string
	}{
		{"正常用户1", 1, "user1"},
		{"正常用户2", 999999, "longusername"},
		{"特殊字符用户", 123, "user_@#$"},
		{"中文用户", 456, "中文用户名"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 生成token
			token, err := GenerateToken(tc.userID, tc.username)
			assert.NoError(t, err)

			// 解析token
			claims, err := ParseToken(token)
			assert.NoError(t, err)
			assert.Equal(t, tc.userID, claims.UserID)
			assert.Equal(t, tc.username, claims.Username)

			// 验证过期时间
			assert.True(t, claims.ExpiresAt.After(time.Now()), "token应该未过期")
		})
	}
}

// TestMultipleTokenGeneration 测试多次生成token
func TestMultipleTokenGeneration(t *testing.T) {
	// 测试不同用户生成的token应该不同
	testCases := []struct {
		userID   int64
		username string
	}{
		{1, "user1"},
		{2, "user2"},
		{3, "user3"},
	}

	tokens := make([]string, len(testCases))
	for i, tc := range testCases {
		token, err := GenerateToken(tc.userID, tc.username)
		assert.NoError(t, err)
		tokens[i] = token
	}

	// 验证所有token都不同（因为用户信息不同）
	for i := 0; i < len(tokens); i++ {
		for j := i + 1; j < len(tokens); j++ {
			assert.NotEqual(t, tokens[i], tokens[j], "不同用户的token应该不同")
		}
	}

	// 验证所有token都能正常解析且信息正确
	for i, token := range tokens {
		claims, err := ParseToken(token)
		assert.NoError(t, err)
		assert.Equal(t, testCases[i].userID, claims.UserID)
		assert.Equal(t, testCases[i].username, claims.Username)
	}
}

// BenchmarkGenerateToken 基准测试：token生成性能
func BenchmarkGenerateToken(b *testing.B) {
	userID := int64(12345)
	username := "testuser"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = GenerateToken(userID, username)
	}
}

// BenchmarkParseToken 基准测试：token解析性能
func BenchmarkParseToken(b *testing.B) {
	userID := int64(12345)
	username := "testuser"
	token, _ := GenerateToken(userID, username)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ParseToken(token)
	}
}
