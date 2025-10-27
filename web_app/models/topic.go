// Package models 定义数据模型
package models

import (
	"time"
)

// Topic 话题模型
type Topic struct {
	ID           int64     `json:"id,string" db:"id"`                // 话题ID（JSON序列化为字符串以避免JavaScript精度丢失）
	UserID       int64     `json:"user_id,string" db:"user_id"`      // 发布者ID
	Username     string    `json:"username" db:"username"`           // 发布者用户名（从users表JOIN）
	Title        string    `json:"title" db:"title"`                 // 话题标题
	Content      string    `json:"content" db:"content"`             // 话题内容
	Category     string    `json:"category" db:"category"`           // 分类（tech/design/discuss/share/product）
	LikeCount    int       `json:"like_count" db:"like_count"`       // 点赞数
	DislikeCount int       `json:"dislike_count" db:"dislike_count"` // 点踩数
	CommentCount int       `json:"comment_count" db:"comment_count"` // 评论数
	ViewCount    int       `json:"view_count" db:"view_count"`       // 浏览数
	CreatedAt    time.Time `json:"created_at" db:"created_at"`       // 创建时间
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`       // 更新时间
}

// CreateTopicRequest 创建话题请求参数
type CreateTopicRequest struct {
	Title    string `json:"title" binding:"required,min=5,max=100"`                              // 标题：5-100个字符
	Content  string `json:"content" binding:"required,min=10"`                                   // 内容：至少10个字符
	Category string `json:"category" binding:"required,oneof=tech design discuss share product"` // 分类
}

// GetTopicsRequest 获取话题列表请求参数
type GetTopicsRequest struct {
	Page     int    `form:"page,default=1"`       // 页码，默认第1页
	PageSize int    `form:"page_size,default=10"` // 每页数量，默认10条
	Sort     string `form:"sort,default=hot"`     // 排序方式：hot/new/like
	Category string `form:"category"`             // 分类筛选（可选）
}

// TopicListResponse 话题列表响应
type TopicListResponse struct {
	Total      int64    `json:"total"`       // 总数
	Page       int      `json:"page"`        // 当前页
	PageSize   int      `json:"page_size"`   // 每页数量
	TotalPages int      `json:"total_pages"` // 总页数
	HasMore    bool     `json:"has_more"`    // 是否有下一页
	Topics     []*Topic `json:"topics"`      // 话题列表
}

// SearchResponse 搜索响应
type SearchResponse struct {
	Total      int64    `json:"total"`       // 总数
	Page       int      `json:"page"`        // 当前页
	PageSize   int      `json:"page_size"`   // 每页数量
	TotalPages int      `json:"total_pages"` // 总页数
	HasMore    bool     `json:"has_more"`    // 是否有下一页
	Topics     []*Topic `json:"topics"`      // 话题列表
	Took       int64    `json:"took"`        // 搜索耗时(毫秒)
}

// TableName 指定表名
func (Topic) TableName() string {
	return "topics"
}
