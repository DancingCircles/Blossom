// Package models 定义数据模型
package models

import (
	"time"
)

// Vote 投票记录模型
type Vote struct {
	ID        int64     `json:"id,string" db:"id"`             // 投票ID（JSON序列化为字符串）
	UserID    int64     `json:"user_id,string" db:"user_id"`   // 用户ID
	TopicID   int64     `json:"topic_id,string" db:"topic_id"` // 话题ID
	VoteType  int       `json:"vote_type" db:"vote_type"`      // 投票类型：1=点赞，-1=点踩
	CreatedAt time.Time `json:"created_at" db:"created_at"`    // 创建时间
}

// VoteRequest 投票请求参数
type VoteRequest struct {
	TopicID  int64 `json:"topic_id" binding:"required"`    // 话题ID
	VoteType int   `json:"vote_type" binding:"oneof=1 -1"` // 投票类型：1或-1
}

// TableName 指定表名
func (Vote) TableName() string {
	return "votes"
}
