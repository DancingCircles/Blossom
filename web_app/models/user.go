// Package models 定义数据模型
package models

import (
	"time"
)

// User 用户模型
type User struct {
	ID        int64     `json:"id,string" db:"id"`          // 用户ID（JSON序列化为字符串以避免JavaScript精度丢失）
	Username  string    `json:"username" db:"username"`     // 用户名
	Email     string    `json:"email" db:"email"`           // 邮箱
	Password  string    `json:"-" db:"password"`            // 密码（不返回给前端）
	CreatedAt time.Time `json:"created_at" db:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"` // 更新时间
}

// RegisterRequest 注册请求参数
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=4,max=20"` // 用户名：4-20个字符
	Email    string `json:"email" binding:"required,email"`           // 邮箱：必须是有效的邮箱格式
	Password string `json:"password" binding:"required,min=6"`        // 密码：至少6个字符
}

// LoginRequest 登录请求参数
type LoginRequest struct {
	Username string `json:"username" binding:"required"` // 用户名
	Password string `json:"password" binding:"required"` // 密码
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token    string `json:"token"`          // JWT token
	Username string `json:"username"`       // 用户名
	UserID   int64  `json:"user_id,string"` // 用户ID（JSON序列化为字符串）
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}
