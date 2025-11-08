// Package tasks 提供定时任务
package tasks

import (
	"context"
	"fmt"
	"time"
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/utils"

	redisv9 "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	// HotRankingKey Redis中存储热榜的key
	HotRankingKey = "hot_ranking"
	// HotRankingSize 热榜保留的话题数量
	HotRankingSize = 1000
)

// UpdateHotRanking 更新热度排行榜
// 每5分钟执行一次，计算最近1000条话题的热度分数
func UpdateHotRanking() error {
	ctx := context.Background()
	startTime := time.Now()

	zap.L().Info("开始更新热度排行榜")

	// 1. 获取最近的话题列表
	topics, err := mysql.GetRecentTopics(HotRankingSize)
	if err != nil {
		zap.L().Error("获取话题列表失败", zap.Error(err))
		return err
	}

	if len(topics) == 0 {
		zap.L().Warn("没有话题需要更新")
		return nil
	}

	// 2. 计算每个话题的热度分数
	rdb := redis.GetClient()
	pipe := rdb.Pipeline()

	for _, topic := range topics {
		score := utils.CalculateHotScore(topic)
		pipe.ZAdd(ctx, HotRankingKey, redisv9.Z{
			Score:  score,
			Member: topic.ID,
		})
	}

	// 3. 批量更新到Redis
	_, err = pipe.Exec(ctx)
	if err != nil {
		zap.L().Error("更新热度排行榜失败", zap.Error(err))
		return err
	}

	// 4. 只保留Top 1000
	err = rdb.ZRemRangeByRank(ctx, HotRankingKey, 0, -HotRankingSize-1).Err()
	if err != nil {
		zap.L().Warn("清理旧数据失败", zap.Error(err))
	}

	// 5. 设置过期时间（1小时）
	err = rdb.Expire(ctx, HotRankingKey, 1*time.Hour).Err()
	if err != nil {
		zap.L().Warn("设置过期时间失败", zap.Error(err))
	}

	duration := time.Since(startTime)
	zap.L().Info("热度排行榜更新完成",
		zap.Int("count", len(topics)),
		zap.Duration("duration", duration))

	return nil
}

// StartHotRankingTask 启动热度排行榜定时任务
func StartHotRankingTask() {
	// 立即执行一次
	if err := UpdateHotRanking(); err != nil {
		zap.L().Error("初始化热度排行榜失败", zap.Error(err))
	}

	// 每5分钟执行一次
	ticker := time.NewTicker(5 * time.Minute)

	go func() {
		for range ticker.C {
			if err := UpdateHotRanking(); err != nil {
				zap.L().Error("更新热度排行榜失败", zap.Error(err))
			}
		}
	}()

	zap.L().Info("热度排行榜定时任务已启动")
}

// GetHotTopicIDs 从Redis获取热门话题ID列表
func GetHotTopicIDs(limit int) ([]int64, error) {
	ctx := context.Background()
	rdb := redis.GetClient()

	// 从Redis ZSet获取热度最高的话题ID
	// ZREVRANGE: 按分数降序获取
	result, err := rdb.ZRevRange(ctx, HotRankingKey, 0, int64(limit-1)).Result()
	if err != nil {
		return nil, err
	}

	// 将字符串ID转换为int64
	ids := make([]int64, 0, len(result))
	for _, idStr := range result {
		var id int64
		_, err := fmt.Sscanf(idStr, "%d", &id)
		if err == nil {
			ids = append(ids, id)
		}
	}

	return ids, nil
}
