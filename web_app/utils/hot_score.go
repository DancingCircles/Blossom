// Package utils 提供工具函数
package utils

import (
	"math"
	"sort"
	"time"
	"web_app/models"
)

// CalculateHotScore 计算话题热度分数
// 基于 Reddit 的热度算法：score = (ups - downs) / (age + 2)^gravity
// ups: 点赞数
// downs: 点踩数
// age: 发布后经过的小时数
// gravity: 重力系数，控制时间衰减速度
func CalculateHotScore(topic *models.Topic) float64 {
	return CalculateHotScoreWithGravity(topic, 1.8)
}

// CalculateHotScoreWithGravity 使用自定义重力系数计算热度
func CalculateHotScoreWithGravity(topic *models.Topic, gravity float64) float64 {
	// 计算投票得分
	votes := topic.LikeCount - topic.DislikeCount

	// 计算发布后经过的小时数
	hoursSinceCreated := time.Since(topic.CreatedAt).Hours()

	// 防止除零错误，至少加2小时
	ageInHours := hoursSinceCreated + 2

	// 计算热度分数
	// 投票数越多，分数越高
	// 时间越久，分数越低（通过幂函数衰减）
	score := float64(votes) / math.Pow(ageInHours, gravity)

	return score
}

// CalculateHackerNewsScore 计算 Hacker News 风格的热度
// score = (votes - 1) / (age + 2)^gravity
func CalculateHackerNewsScore(topic *models.Topic) float64 {
	votes := float64(topic.LikeCount - 1)
	if votes < 0 {
		votes = 0
	}

	hoursSinceCreated := time.Since(topic.CreatedAt).Hours()
	ageInHours := hoursSinceCreated + 2

	return votes / math.Pow(ageInHours, 1.8)
}

// CalculateWilsonScore 计算 Wilson 置信区间分数
// 用于处理投票数较少的情况，更加平衡
func CalculateWilsonScore(likes, dislikes int) float64 {
	if likes+dislikes == 0 {
		return 0
	}

	n := float64(likes + dislikes)
	p := float64(likes) / n

	// z = 1.96 for 95% confidence
	z := 1.96

	// Wilson score interval
	score := (p + z*z/(2*n) - z*math.Sqrt((p*(1-p)+z*z/(4*n))/n)) / (1 + z*z/n)

	return score
}

// CalculateEngagementScore 计算参与度分数
// 综合考虑点赞、评论、浏览量
func CalculateEngagementScore(topic *models.Topic) float64 {
	// 权重配置
	likeWeight := 1.0
	commentWeight := 2.0 // 评论权重更高
	viewWeight := 0.01   // 浏览权重较低

	// 计算参与度
	engagement := float64(topic.LikeCount)*likeWeight +
		float64(topic.CommentCount)*commentWeight +
		float64(topic.ViewCount)*viewWeight

	// 时间衰减
	hoursSinceCreated := time.Since(topic.CreatedAt).Hours()
	ageInHours := hoursSinceCreated + 2

	return engagement / math.Pow(ageInHours, 1.5)
}

// RankTopics 对话题列表按热度排序
func RankTopics(topics []*models.Topic) []*models.Topic {
	// 创建副本避免修改原数组
	ranked := make([]*models.Topic, len(topics))
	copy(ranked, topics)

	// 预计算所有话题的热度分数（避免重复计算）
	scores := make(map[int64]float64, len(ranked))
	for _, topic := range ranked {
		scores[topic.ID] = CalculateHotScore(topic)
	}

	// 使用快速排序（sort.Slice），时间复杂度 O(n log n)
	sort.Slice(ranked, func(i, j int) bool {
		return scores[ranked[i].ID] > scores[ranked[j].ID] // 降序排列
	})

	return ranked
}
