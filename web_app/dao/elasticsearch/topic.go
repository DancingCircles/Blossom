// Package elasticsearch 提供话题索引操作
package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"web_app/models"

	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
)

// TopicDocument ES中的话题文档结构
type TopicDocument struct {
	TopicID      string `json:"topic_id"`
	UserID       string `json:"user_id"`
	Title        string `json:"title"`
	Content      string `json:"content"`
	Category     string `json:"category"`
	CreatedAt    string `json:"created_at"` // 使用string以匹配ES的日期格式
	UpdatedAt    string `json:"updated_at"` // 使用string以匹配ES的日期格式
	ViewCount    int    `json:"view_count"`
	CommentCount int    `json:"comment_count"`
}

// IndexTopic 将话题添加到ES索引
func IndexTopic(topic *models.Topic) error {
	ctx := context.Background()

	// 时间格式（匹配ES mapping）
	timeFormat := "2006-01-02 15:04:05"

	// 构建文档
	doc := TopicDocument{
		TopicID:      fmt.Sprintf("%d", topic.ID),
		UserID:       fmt.Sprintf("%d", topic.UserID),
		Title:        topic.Title,
		Content:      topic.Content,
		Category:     topic.Category,
		CreatedAt:    topic.CreatedAt.Format(timeFormat),
		UpdatedAt:    topic.UpdatedAt.Format(timeFormat),
		ViewCount:    topic.ViewCount,
		CommentCount: topic.CommentCount,
	}

	// 索引文档
	_, err := client.Index().
		Index(index).
		Id(doc.TopicID).
		BodyJson(doc).
		Refresh("true"). // 立即刷新，使其可搜索
		Do(ctx)

	if err != nil {
		zap.L().Error("索引话题失败", zap.Error(err), zap.String("topic_id", doc.TopicID))
		return err
	}

	zap.L().Debug("话题已索引", zap.String("topic_id", doc.TopicID))
	return nil
}

// UpdateTopic 更新ES中的话题
func UpdateTopic(topic *models.Topic) error {
	ctx := context.Background()

	// 时间格式（匹配ES mapping）
	timeFormat := "2006-01-02 15:04:05"

	// 构建更新文档
	doc := map[string]interface{}{
		"title":         topic.Title,
		"content":       topic.Content,
		"category":      topic.Category,
		"updated_at":    topic.UpdatedAt.Format(timeFormat),
		"view_count":    topic.ViewCount,
		"comment_count": topic.CommentCount,
	}

	topicID := fmt.Sprintf("%d", topic.ID)

	// 更新文档
	_, err := client.Update().
		Index(index).
		Id(topicID).
		Doc(doc).
		Refresh("true").
		Do(ctx)

	if err != nil {
		zap.L().Error("更新话题索引失败", zap.Error(err), zap.String("topic_id", topicID))
		return err
	}

	zap.L().Debug("话题索引已更新", zap.String("topic_id", topicID))
	return nil
}

// DeleteTopic 从ES中删除话题
func DeleteTopic(topicID int64) error {
	ctx := context.Background()

	id := fmt.Sprintf("%d", topicID)

	// 删除文档
	_, err := client.Delete().
		Index(index).
		Id(id).
		Refresh("true").
		Do(ctx)

	if err != nil {
		// 如果文档不存在，不算错误
		if elastic.IsNotFound(err) {
			zap.L().Warn("话题索引不存在", zap.String("topic_id", id))
			return nil
		}
		zap.L().Error("删除话题索引失败", zap.Error(err), zap.String("topic_id", id))
		return err
	}

	zap.L().Debug("话题索引已删除", zap.String("topic_id", id))
	return nil
}

// BulkIndexTopics 批量索引话题
func BulkIndexTopics(topics []*models.Topic) error {
	if len(topics) == 0 {
		return nil
	}

	ctx := context.Background()
	bulkRequest := client.Bulk()

	// 时间格式（匹配ES mapping）
	timeFormat := "2006-01-02 15:04:05"

	for _, topic := range topics {
		doc := TopicDocument{
			TopicID:      fmt.Sprintf("%d", topic.ID),
			UserID:       fmt.Sprintf("%d", topic.UserID),
			Title:        topic.Title,
			Content:      topic.Content,
			Category:     topic.Category,
			CreatedAt:    topic.CreatedAt.Format(timeFormat),
			UpdatedAt:    topic.UpdatedAt.Format(timeFormat),
			ViewCount:    topic.ViewCount,
			CommentCount: topic.CommentCount,
		}

		req := elastic.NewBulkIndexRequest().
			Index(index).
			Id(doc.TopicID).
			Doc(doc)

		bulkRequest = bulkRequest.Add(req)
	}

	// 执行批量操作
	bulkResponse, err := bulkRequest.Refresh("true").Do(ctx)
	if err != nil {
		zap.L().Error("批量索引话题失败", zap.Error(err))
		return err
	}

	// 检查是否有失败的操作
	if bulkResponse.Errors {
		zap.L().Warn("部分话题索引失败", zap.Int("failed", len(bulkResponse.Failed())))
		for _, item := range bulkResponse.Failed() {
			zap.L().Error("索引失败",
				zap.String("id", item.Id),
				zap.String("error", item.Error.Reason))
		}
	}

	zap.L().Info("批量索引完成",
		zap.Int("total", len(topics)),
		zap.Int("success", len(bulkResponse.Succeeded())),
		zap.Int("failed", len(bulkResponse.Failed())))

	return nil
}

// UpdateTopicCommentCount 更新话题的评论数
func UpdateTopicCommentCount(topicID int64, commentCount int) error {
	ctx := context.Background()

	id := fmt.Sprintf("%d", topicID)

	// 只更新评论数字段
	doc := map[string]interface{}{
		"comment_count": commentCount,
	}

	// 更新文档
	_, err := client.Update().
		Index(index).
		Id(id).
		Doc(doc).
		Refresh("true").
		Do(ctx)

	if err != nil {
		// 如果文档不存在，记录警告但不返回错误
		if elastic.IsNotFound(err) {
			zap.L().Warn("话题索引不存在，跳过更新评论数",
				zap.String("topic_id", id),
				zap.Int("comment_count", commentCount))
			return nil
		}
		zap.L().Error("更新话题评论数失败",
			zap.Error(err),
			zap.String("topic_id", id),
			zap.Int("comment_count", commentCount))
		return err
	}

	zap.L().Debug("话题评论数已更新",
		zap.String("topic_id", id),
		zap.Int("comment_count", commentCount))
	return nil
}

// GetTopicByID 从ES中获取话题（用于调试）
func GetTopicByID(topicID int64) (*TopicDocument, error) {
	ctx := context.Background()

	id := fmt.Sprintf("%d", topicID)

	result, err := client.Get().
		Index(index).
		Id(id).
		Do(ctx)

	if err != nil {
		if elastic.IsNotFound(err) {
			return nil, fmt.Errorf("话题不存在")
		}
		return nil, err
	}

	var doc TopicDocument
	err = json.Unmarshal(result.Source, &doc)
	if err != nil {
		return nil, err
	}

	return &doc, nil
}
