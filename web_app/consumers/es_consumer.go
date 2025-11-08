// Package consumers 提供Kafka消费者实现
package consumers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"web_app/dao/elasticsearch"
	"web_app/dao/mysql"
	"web_app/models"

	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// CanalMessage Canal消息格式
type CanalMessage struct {
	Type      string                   `json:"type"`      // INSERT/UPDATE/DELETE
	Database  string                   `json:"database"`  // 数据库名
	Table     string                   `json:"table"`     // 表名
	Data      []map[string]interface{} `json:"data"`      // 新数据
	Old       []map[string]interface{} `json:"old"`       // 旧数据（UPDATE时有值）
	IsDdl     bool                     `json:"isDdl"`     // 是否DDL操作
	Es        int64                    `json:"es"`        // 执行时间戳（毫秒）
	Ts        int64                    `json:"ts"`        // Canal接收时间戳（毫秒）
	MysqlType map[string]string        `json:"mysqlType"` // MySQL字段类型
	SqlType   map[string]int           `json:"sqlType"`   // SQL字段类型
}

// ESConsumer ES同步消费者
type ESConsumer struct {
	reader *kafka.Reader
}

// NewESConsumer 创建ES同步消费者
func NewESConsumer() *ESConsumer {
	brokers := viper.GetStringSlice("kafka.brokers")
	topic := viper.GetString("kafka.topic")
	groupID := viper.GetString("kafka.group_id")

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        groupID,
		MinBytes:       1,
		MaxBytes:       10e6, // 10MB
		CommitInterval: time.Second,
		StartOffset:    kafka.LastOffset,
	})

	return &ESConsumer{
		reader: reader,
	}
}

// Start 启动消费者
func (c *ESConsumer) Start() {
	zap.L().Info("ES同步消费者已启动",
		zap.Strings("brokers", viper.GetStringSlice("kafka.brokers")),
		zap.String("topic", viper.GetString("kafka.topic")),
		zap.String("group_id", viper.GetString("kafka.group_id")))

	for {
		msg, err := c.reader.ReadMessage(context.Background())
		if err != nil {
			zap.L().Error("读取Kafka消息失败", zap.Error(err))
			time.Sleep(time.Second)
			continue
		}

		// 处理消息
		if err := c.processMessage(msg.Value); err != nil {
			zap.L().Error("处理消息失败",
				zap.Error(err),
				zap.String("message", string(msg.Value)))
		}
	}
}

// processMessage 处理Canal消息
func (c *ESConsumer) processMessage(data []byte) error {
	var canalMsg CanalMessage
	if err := json.Unmarshal(data, &canalMsg); err != nil {
		return fmt.Errorf("解析Canal消息失败: %w", err)
	}

	// 跳过DDL操作
	if canalMsg.IsDdl {
		return nil
	}

	// 只处理topics和comments表
	switch canalMsg.Table {
	case "topics":
		return c.handleTopicChange(&canalMsg)
	case "comments":
		return c.handleCommentChange(&canalMsg)
	default:
		// 忽略其他表
		return nil
	}
}

// handleTopicChange 处理话题变更
func (c *ESConsumer) handleTopicChange(msg *CanalMessage) error {
	switch msg.Type {
	case "INSERT", "UPDATE":
		// 新增或更新话题
		for _, data := range msg.Data {
			topic, err := c.parseTopicFromData(data)
			if err != nil {
				zap.L().Error("解析话题数据失败", zap.Error(err))
				continue
			}

			// 同步到ES
			if err := elasticsearch.IndexTopic(topic); err != nil {
				zap.L().Error("同步话题到ES失败",
					zap.Error(err),
					zap.Int64("topic_id", topic.ID))
				return err
			}

			zap.L().Info("话题已同步到ES",
				zap.Int64("topic_id", topic.ID),
				zap.String("title", topic.Title),
				zap.String("operation", msg.Type))
		}

	case "DELETE":
		// 删除话题
		for _, data := range msg.Data {
			topicID := c.parseInt64(data["id"])
			if err := elasticsearch.DeleteTopic(topicID); err != nil {
				zap.L().Error("从ES删除话题失败",
					zap.Error(err),
					zap.Int64("topic_id", topicID))
				return err
			}

			zap.L().Info("话题已从ES删除",
				zap.Int64("topic_id", topicID))
		}
	}

	return nil
}

// handleCommentChange 处理评论变更（更新话题的评论数）
func (c *ESConsumer) handleCommentChange(msg *CanalMessage) error {
	zap.L().Debug("检测到评论变更",
		zap.String("operation", msg.Type),
		zap.Int("count", len(msg.Data)))

	// 收集需要更新的话题ID（使用map去重）
	topicIDsMap := make(map[int64]bool)

	// 从所有变更的评论中提取topic_id
	for _, data := range msg.Data {
		topicID := c.parseInt64(data["topic_id"])
		if topicID > 0 {
			topicIDsMap[topicID] = true
		}
	}

	// 如果没有有效的topic_id，直接返回
	if len(topicIDsMap) == 0 {
		zap.L().Debug("没有需要更新的话题")
		return nil
	}

	// 对每个话题，查询最新的评论数并更新到ES
	for topicID := range topicIDsMap {
		// 1. 从MySQL查询该话题的最新评论数
		count, err := mysql.CountCommentsByTopicID(topicID)
		if err != nil {
			zap.L().Error("查询话题评论数失败",
				zap.Error(err),
				zap.Int64("topic_id", topicID))
			continue // 继续处理其他话题
		}

		// 2. 更新ES中的话题评论数
		err = elasticsearch.UpdateTopicCommentCount(topicID, int(count))
		if err != nil {
			zap.L().Error("更新ES话题评论数失败",
				zap.Error(err),
				zap.Int64("topic_id", topicID),
				zap.Int64("comment_count", count))
			continue
		}

		zap.L().Info("话题评论数已同步到ES",
			zap.Int64("topic_id", topicID),
			zap.Int64("comment_count", count),
			zap.String("operation", msg.Type))
	}

	return nil
}

// parseTopicFromData 从Canal数据解析Topic结构
func (c *ESConsumer) parseTopicFromData(data map[string]interface{}) (*models.Topic, error) {
	topic := &models.Topic{
		ID:           c.parseInt64(data["id"]),
		UserID:       c.parseInt64(data["user_id"]),
		Title:        c.parseString(data["title"]),
		Content:      c.parseString(data["content"]),
		Category:     c.parseString(data["category"]),
		ViewCount:    c.parseInt(data["view_count"]),
		CommentCount: c.parseInt(data["comment_count"]),
	}

	// 解析时间字段
	if createdAt := c.parseString(data["created_at"]); createdAt != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", createdAt); err == nil {
			topic.CreatedAt = t
		}
	}

	if updatedAt := c.parseString(data["updated_at"]); updatedAt != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", updatedAt); err == nil {
			topic.UpdatedAt = t
		}
	}

	return topic, nil
}

// parseInt64 辅助函数：安全地将interface{}转换为int64
func (c *ESConsumer) parseInt64(val interface{}) int64 {
	switch v := val.(type) {
	case float64:
		return int64(v)
	case int64:
		return v
	case int:
		return int64(v)
	case string:
		var result int64
		fmt.Sscanf(v, "%d", &result)
		return result
	default:
		return 0
	}
}

// parseInt 辅助函数：安全地将interface{}转换为int
func (c *ESConsumer) parseInt(val interface{}) int {
	return int(c.parseInt64(val))
}

// parseString 辅助函数：安全地将interface{}转换为string
func (c *ESConsumer) parseString(val interface{}) string {
	if val == nil {
		return ""
	}
	if str, ok := val.(string); ok {
		return str
	}
	return fmt.Sprintf("%v", val)
}

// Close 关闭消费者
func (c *ESConsumer) Close() error {
	if c.reader != nil {
		return c.reader.Close()
	}
	return nil
}

// StartESConsumer 启动ES同步消费者（供main.go调用）
func StartESConsumer() {
	consumer := NewESConsumer()
	consumer.Start()
}
