// Package logic 提供搜索业务逻辑
package logic

import (
	"fmt"
	"strconv"
	"time"
	"web_app/dao/elasticsearch"
	"web_app/models"
)

// SearchTopics 搜索话题
func SearchTopics(keyword, category string, page, pageSize int, sortBy string) (*models.SearchResponse, error) {
	// 参数验证
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 调用ES搜索
	req := &elasticsearch.SearchRequest{
		Keyword:  keyword,
		Category: category,
		Page:     page,
		PageSize: pageSize,
		SortBy:   sortBy,
	}

	esResp, err := elasticsearch.SearchTopics(req)
	if err != nil {
		return nil, err
	}

	// 转换为models.Topic格式
	topics := make([]*models.Topic, 0, len(esResp.Topics))
	timeFormat := "2006-01-02 15:04:05"

	for _, doc := range esResp.Topics {
		topicID, _ := strconv.ParseInt(doc.TopicID, 10, 64)
		userID, _ := strconv.ParseInt(doc.UserID, 10, 64)

		// 解析时间字符串
		createdAt, _ := time.Parse(timeFormat, doc.CreatedAt)
		updatedAt, _ := time.Parse(timeFormat, doc.UpdatedAt)

		topic := &models.Topic{
			ID:           topicID,
			UserID:       userID,
			Title:        doc.Title,
			Content:      doc.Content,
			Category:     doc.Category,
			CreatedAt:    createdAt,
			UpdatedAt:    updatedAt,
			ViewCount:    doc.ViewCount,
			CommentCount: doc.CommentCount,
		}
		topics = append(topics, topic)
	}

	// 构建响应
	totalPages := int((esResp.Total + int64(pageSize) - 1) / int64(pageSize))
	hasMore := page < totalPages

	response := &models.SearchResponse{
		Total:      esResp.Total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		HasMore:    hasMore,
		Topics:     topics,
		Took:       esResp.Took,
	}

	return response, nil
}

// SuggestTopics 搜索建议
func SuggestTopics(prefix string) ([]string, error) {
	if prefix == "" {
		return []string{}, nil
	}

	suggestions, err := elasticsearch.SuggestTopics(prefix, 10)
	if err != nil {
		return nil, fmt.Errorf("获取搜索建议失败: %w", err)
	}

	return suggestions, nil
}

// GetHotTopicsByCategory 获取分类热门话题
func GetHotTopicsByCategory(category string, size int) ([]*models.Topic, error) {
	if size <= 0 || size > 50 {
		size = 10
	}

	docs, err := elasticsearch.GetTopicsByCategory(category, size)
	if err != nil {
		return nil, fmt.Errorf("获取分类热门话题失败: %w", err)
	}

	topics := make([]*models.Topic, 0, len(docs))
	timeFormat := "2006-01-02 15:04:05"

	for _, doc := range docs {
		topicID, _ := strconv.ParseInt(doc.TopicID, 10, 64)
		userID, _ := strconv.ParseInt(doc.UserID, 10, 64)

		// 解析时间字符串
		createdAt, _ := time.Parse(timeFormat, doc.CreatedAt)
		updatedAt, _ := time.Parse(timeFormat, doc.UpdatedAt)

		topic := &models.Topic{
			ID:           topicID,
			UserID:       userID,
			Title:        doc.Title,
			Content:      doc.Content,
			Category:     doc.Category,
			CreatedAt:    createdAt,
			UpdatedAt:    updatedAt,
			ViewCount:    doc.ViewCount,
			CommentCount: doc.CommentCount,
		}
		topics = append(topics, topic)
	}

	return topics, nil
}

// GetCategoryStats 获取分类统计
func GetCategoryStats() (map[string]int64, error) {
	stats, err := elasticsearch.CountByCategory()
	if err != nil {
		return nil, fmt.Errorf("获取分类统计失败: %w", err)
	}

	return stats, nil
}
