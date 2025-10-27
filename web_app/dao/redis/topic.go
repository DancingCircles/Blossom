package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"web_app/models"
)

const (
	// 缓存键前缀
	topicDetailPrefix = "topic:detail:"  // 话题详情缓存键前缀
	topicListPrefix   = "topic:list:"    // 话题列表缓存键前缀
	hotTopicsKey      = "topic:list:hot" // 热门话题列表缓存键

	// 缓存过期时间
	topicDetailTTL = 10 * time.Minute // 话题详情缓存10分钟
	topicListTTL   = 5 * time.Minute  // 话题列表缓存5分钟
)

// CacheTopicDetail 缓存话题详情
func CacheTopicDetail(topic *models.Topic) error {
	ctx := context.Background()
	key := fmt.Sprintf("%s%d", topicDetailPrefix, topic.ID)

	// 序列化话题数据
	data, err := json.Marshal(topic)
	if err != nil {
		return err
	}

	// 写入Redis
	return rdb.Set(ctx, key, data, topicDetailTTL).Err()
}

// GetTopicDetailCache 获取话题详情缓存
func GetTopicDetailCache(topicID int64) (*models.Topic, error) {
	ctx := context.Background()
	key := fmt.Sprintf("%s%d", topicDetailPrefix, topicID)

	// 从Redis读取
	data, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, err // 包括redis.Nil（缓存未命中）
	}

	// 反序列化
	var topic models.Topic
	if err := json.Unmarshal([]byte(data), &topic); err != nil {
		return nil, err
	}

	return &topic, nil
}

// DeleteTopicCache 删除话题详情缓存
func DeleteTopicCache(topicID int64) error {
	ctx := context.Background()
	key := fmt.Sprintf("%s%d", topicDetailPrefix, topicID)
	return rdb.Del(ctx, key).Err()
}

// CacheTopicList 缓存话题列表
func CacheTopicList(cacheKey string, topics []*models.Topic, total int64) error {
	ctx := context.Background()

	// 构建缓存数据结构
	cacheData := map[string]interface{}{
		"topics": topics,
		"total":  total,
	}

	// 序列化
	data, err := json.Marshal(cacheData)
	if err != nil {
		return err
	}

	// 写入Redis
	return rdb.Set(ctx, cacheKey, data, topicListTTL).Err()
}

// GetTopicListCache 获取话题列表缓存
func GetTopicListCache(cacheKey string) ([]*models.Topic, int64, error) {
	ctx := context.Background()

	// 从Redis读取
	data, err := rdb.Get(ctx, cacheKey).Result()
	if err != nil {
		return nil, 0, err
	}

	// 反序列化
	var cacheData struct {
		Topics []*models.Topic `json:"topics"`
		Total  int64           `json:"total"`
	}
	if err := json.Unmarshal([]byte(data), &cacheData); err != nil {
		return nil, 0, err
	}

	return cacheData.Topics, cacheData.Total, nil
}

// DeleteTopicListCache 删除话题列表缓存（支持通配符）
func DeleteTopicListCache(pattern string) error {
	ctx := context.Background()

	// 如果没有通配符，直接删除
	if pattern == "" {
		pattern = topicListPrefix + "*"
	}

	// 扫描匹配的键
	iter := rdb.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := rdb.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}

	return iter.Err()
}

// BuildTopicListCacheKey 构建话题列表缓存键
func BuildTopicListCacheKey(req *models.GetTopicsRequest) string {
	return fmt.Sprintf("%spage:%d:size:%d:sort:%s:category:%s",
		topicListPrefix, req.Page, req.PageSize, req.Sort, req.Category)
}

// DeleteAllTopicListCache 删除所有话题列表缓存
func DeleteAllTopicListCache() error {
	return DeleteTopicListCache(topicListPrefix + "*")
}
