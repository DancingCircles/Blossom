// Package utils_test 提供工具函数测试
package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestHashPassword 测试密码加密
func TestHashPassword(t *testing.T) {
	password := "testPassword123"

	// 测试正常加密
	hashedPassword, err := HashPassword(password)
	assert.NoError(t, err, "密码加密应该成功")
	assert.NotEmpty(t, hashedPassword, "加密后的密码不应该为空")
	assert.NotEqual(t, password, hashedPassword, "加密后的密码应该与原密码不同")

	// 测试相同密码加密两次结果不同（salt机制）
	hashedPassword2, err := HashPassword(password)
	assert.NoError(t, err)
	assert.NotEqual(t, hashedPassword, hashedPassword2, "相同密码两次加密结果应该不同")
}

// TestCheckPassword 测试密码验证
func TestCheckPassword(t *testing.T) {
	password := "testPassword123"
	wrongPassword := "wrongPassword"

	// 先加密密码
	hashedPassword, err := HashPassword(password)
	assert.NoError(t, err)

	// 测试正确密码验证
	err = CheckPassword(hashedPassword, password)
	assert.NoError(t, err, "正确密码验证应该成功")

	// 测试错误密码验证
	err = CheckPassword(hashedPassword, wrongPassword)
	assert.Error(t, err, "错误密码验证应该失败")
}

// TestHashPasswordEmpty 测试空密码
func TestHashPasswordEmpty(t *testing.T) {
	// 测试空密码（bcrypt允许空密码）
	hashedPassword, err := HashPassword("")
	assert.NoError(t, err, "空密码加密应该成功")
	assert.NotEmpty(t, hashedPassword, "加密后的密码不应该为空")
}

// TestCheckPasswordWithInvalidHash 测试无效的哈希值
func TestCheckPasswordWithInvalidHash(t *testing.T) {
	password := "testPassword123"
	invalidHash := "invalid_hash_value"

	// 测试无效哈希值
	err := CheckPassword(invalidHash, password)
	assert.Error(t, err, "无效哈希值验证应该失败")
}

// BenchmarkHashPassword 基准测试：密码加密性能
func BenchmarkHashPassword(b *testing.B) {
	password := "testPassword123"
	for i := 0; i < b.N; i++ {
		_, _ = HashPassword(password)
	}
}

// BenchmarkCheckPassword 基准测试：密码验证性能
func BenchmarkCheckPassword(b *testing.B) {
	password := "testPassword123"
	hashedPassword, _ := HashPassword(password)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = CheckPassword(hashedPassword, password)
	}
}
