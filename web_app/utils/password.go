// Package utils 提供工具函数
package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 加密密码
// 使用bcrypt算法加密密码
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPassword 验证密码
// 比较明文密码和加密后的密码是否匹配
func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
