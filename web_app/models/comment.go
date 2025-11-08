// Package models 定义数据模型
package models

import (
	"time"
)

// Comment 评论模型
type Comment struct {
	ID        int64     `json:"id,string" db:"id"`               // 评论ID（JSON序列化为字符串以避免JavaScript精度丢失）
	TopicID   int64     `json:"topic_id,string" db:"topic_id"`   // 话题ID
	UserID    int64     `json:"user_id,string" db:"user_id"`     // 用户ID
	Username  string    `json:"username" db:"username"`          // 用户名（从users表JOIN）
	Content   string    `json:"content" db:"content"`            // 评论内容
	ParentID  *int64    `json:"parent_id,string" db:"parent_id"` // 父评论ID（用于回复，可为空）
	CreatedAt time.Time `json:"created_at" db:"created_at"`      // 创建时间
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`      // 更新时间
}

// CreateCommentRequest 创建评论请求参数
type CreateCommentRequest struct {
	Content  string `json:"content" binding:"required,min=1,max=1000"` // 评论内容：1-1000个字符
	ParentID *int64 `json:"parent_id,string,omitempty"`                // 父评论ID（可选，用于回复）
}

// CommentListResponse 评论列表响应
type CommentListResponse struct {
	Total      int64      `json:"total"`       // 总数
	Page       int        `json:"page"`        // 当前页
	PageSize   int        `json:"page_size"`   // 每页数量
	TotalPages int        `json:"total_pages"` // 总页数
	HasMore    bool       `json:"has_more"`    // 是否有下一页
	Comments   []*Comment `json:"comments"`    // 评论列表
}

// GetCommentsRequest 获取评论列表请求参数
type GetCommentsRequest struct {
	Page     int `form:"page,default=1"`       // 页码，默认第1页
	PageSize int `form:"page_size,default=20"` // 每页数量，默认20条
}

// TableName 指定表名
func (Comment) TableName() string {
	return "comments"
}
