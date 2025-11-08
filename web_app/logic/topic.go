package logic

import (
	"context"
	"errors"
	"fmt"
	"time"
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/models"
	"web_app/tasks"
	"web_app/utils"

	"go.uber.org/zap"
)

// CreateTopic 创建话题
func CreateTopic(userID int64, req *models.CreateTopicRequest) error {
	// 生成雪花算法ID
	topicID := utils.GenerateID()

	// 调用dao层插入话题（显式设置创建时间为当前时间）
	now := time.Now()
	topic := &models.Topic{
		ID:        topicID,
		UserID:    userID,
		Title:     req.Title,
		Content:   req.Content,
		Category:  req.Category,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := mysql.InsertTopic(topic)
	if err != nil {
		zap.L().Error("插入话题失败", zap.Error(err))
		return errors.New("插入话题失败")
	}

	// 异步清除话题列表缓存
	go func() {
		if err := redis.DeleteAllTopicListCache(); err != nil {
			zap.L().Warn("清除话题列表缓存失败", zap.Error(err))
		}
	}()

	// ES同步由Canal+Kafka自动处理，无需手动同步

	return nil
}

// GetTopics 获取话题列表
func GetTopics(req *models.GetTopicsRequest) ([]*models.Topic, int64, error) {
	// 1. 尝试从Redis缓存获取
	cacheKey := redis.BuildTopicListCacheKey(req)
	topicList, total, err := redis.GetTopicListCache(cacheKey)
	if err == nil {
		// 缓存命中
		zap.L().Debug("话题列表缓存命中", zap.String("cache_key", cacheKey))
		return topicList, total, nil
	}

	// 2. 缓存未命中，从数据库查询
	topicList, total, err = mysql.GetTopics(req)
	if err != nil {
		zap.L().Error("查询话题失败", zap.Error(err))
		return nil, 0, errors.New("查询话题失败")
	}

	// 3. 异步写入缓存
	go func() {
		if err := redis.CacheTopicList(cacheKey, topicList, total); err != nil {
			zap.L().Warn("缓存话题列表失败", zap.Error(err))
		}
	}()

	return topicList, total, nil
}

// GetTopicByID 获取话题详情
func GetTopicByID(topicID int64) (*models.Topic, error) {
	// 1. 尝试从Redis缓存获取
	topic, err := redis.GetTopicDetailCache(topicID)
	if err == nil {
		// 缓存命中
		zap.L().Debug("话题详情缓存命中", zap.Int64("topic_id", topicID))
		// 即使缓存命中，也要增加浏览数
		go func() {
			if err := mysql.UpdateTopicViewCount(topicID); err != nil {
				zap.L().Warn("更新浏览数失败", zap.Error(err))
			}
		}()
		return topic, nil
	}

	// 2. 缓存未命中，从数据库查询
	topic, err = mysql.GetTopicByID(topicID)
	if err != nil {
		// 如果是话题不存在，使用WARN级别；其他错误使用ERROR级别
		if err.Error() == "话题不存在" {
			zap.L().Warn("查询话题详情失败", zap.Int64("topic_id", topicID), zap.Error(err))
		} else {
			zap.L().Error("查询话题详情失败", zap.Error(err))
		}
		return nil, err
	}

	// 3. 异步写入缓存
	go func() {
		if err := redis.CacheTopicDetail(topic); err != nil {
			zap.L().Warn("缓存话题详情失败", zap.Error(err))
		}
	}()

	// 4. 增加浏览数（异步处理，不影响主流程）
	go func() {
		if err := mysql.UpdateTopicViewCount(topicID); err != nil {
			zap.L().Warn("更新浏览数失败", zap.Error(err))
		}
	}()

	return topic, nil
}

// VoteTopic 投票（点赞/点踩）
func VoteTopic(userID, topicID int64, voteType string) error {
	// 使用分布式锁防止并发投票导致的数据不一致
	ctx := context.Background()
	lockKey := fmt.Sprintf("lock:vote:%d:%d", topicID, userID)

	// 使用 WithLock 自动管理锁的获取和释放
	return utils.WithLock(ctx, redis.GetClient(), lockKey, 3*time.Second, func() error {
		// 1. 验证话题是否存在
		_, err := mysql.GetTopicByID(topicID)
		if err != nil {
			return errors.New("话题不存在")
		}

		// 2. 转换投票类型
		var voteValue int
		if voteType == "like" {
			voteValue = 1
		} else if voteType == "dislike" {
			voteValue = -1
		} else {
			return errors.New("无效的投票类型")
		}

		// 3. 查询用户是否已投票
		existingVote, err := mysql.GetUserVote(userID, topicID)
		if err != nil {
			zap.L().Error("查询投票记录失败", zap.Error(err))
			return errors.New("查询投票记录失败")
		}

		// 4. 处理投票逻辑
		return processVoteLogic(userID, topicID, voteValue, existingVote)
	})
}

// processVoteLogic 处理具体的投票逻辑（内部函数）
func processVoteLogic(userID, topicID int64, voteValue int, existingVote *models.Vote) error {
	if existingVote == nil {
		// 未投票，新增投票
		voteID := utils.GenerateID()
		if err := mysql.InsertVote(&models.Vote{
			ID:       voteID,
			UserID:   userID,
			TopicID:  topicID,
			VoteType: voteValue,
		}); err != nil {
			return errors.New("投票失败")
		}

		// 更新话题计数
		if voteValue == 1 {
			if err := mysql.UpdateTopicLikeCount(topicID, 1); err != nil {
				zap.L().Error("更新点赞数失败", zap.Error(err))
			}
		} else {
			if err := mysql.UpdateTopicDislikeCount(topicID, 1); err != nil {
				zap.L().Error("更新点踩数失败", zap.Error(err))
			}
		}
	} else if existingVote.VoteType == voteValue {
		// 已投相同类型的票，取消投票
		if err := mysql.DeleteVote(userID, topicID); err != nil {
			return errors.New("取消投票失败")
		}

		// 更新话题计数
		if voteValue == 1 {
			if err := mysql.UpdateTopicLikeCount(topicID, -1); err != nil {
				zap.L().Error("更新点赞数失败", zap.Error(err))
			}
		} else {
			if err := mysql.UpdateTopicDislikeCount(topicID, -1); err != nil {
				zap.L().Error("更新点踩数失败", zap.Error(err))
			}
		}
	} else {
		// 已投不同类型的票，更新投票
		if err := mysql.UpdateVote(userID, topicID, voteValue); err != nil {
			return errors.New("更新投票失败")
		}

		// 更新话题计数（减少原来的，增加新的）
		if existingVote.VoteType == 1 {
			if err := mysql.UpdateTopicLikeCount(topicID, -1); err != nil {
				zap.L().Error("更新点赞数失败", zap.Error(err))
			}
			if err := mysql.UpdateTopicDislikeCount(topicID, 1); err != nil {
				zap.L().Error("更新点踩数失败", zap.Error(err))
			}
		} else {
			if err := mysql.UpdateTopicDislikeCount(topicID, -1); err != nil {
				zap.L().Error("更新点踩数失败", zap.Error(err))
			}
			if err := mysql.UpdateTopicLikeCount(topicID, 1); err != nil {
				zap.L().Error("更新点赞数失败", zap.Error(err))
			}
		}
	}

	// 异步清除该话题的缓存和列表缓存
	go func() {
		// 清除话题详情缓存
		if err := redis.DeleteTopicCache(topicID); err != nil {
			zap.L().Warn("清除话题详情缓存失败", zap.Error(err))
		}
		// 清除话题列表缓存
		if err := redis.DeleteAllTopicListCache(); err != nil {
			zap.L().Warn("清除话题列表缓存失败", zap.Error(err))
		}
	}()

	return nil
}

// GetHotTopics 获取热门话题列表
func GetHotTopics(limit int) ([]*models.Topic, error) {
	// 参数验证
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	// 从Redis获取热门话题ID列表
	ids, err := tasks.GetHotTopicIDs(limit)
	if err != nil || len(ids) == 0 {
		// 如果Redis中没有数据，回退到从数据库获取（限制数量避免性能问题）
		zap.L().Warn("从Redis获取热榜失败，回退到数据库查询", zap.Error(err))
		fallbackLimit := limit * 3 // 获取3倍数量用于排序，避免查询过多
		if fallbackLimit > 100 {
			fallbackLimit = 100
		}
		topics, err := mysql.GetRecentTopics(fallbackLimit)
		if err != nil {
			return nil, err
		}
		// 按热度排序并返回前limit条
		ranked := utils.RankTopics(topics)
		if len(ranked) > limit {
			return ranked[:limit], nil
		}
		return ranked, nil
	}

	// 根据ID列表从数据库获取话题详情
	topics, err := mysql.GetTopicsByIDs(ids)
	if err != nil {
		return nil, err
	}

	// 保持Redis中的顺序
	topicMap := make(map[int64]*models.Topic)
	for _, topic := range topics {
		topicMap[topic.ID] = topic
	}

	orderedTopics := make([]*models.Topic, 0, len(ids))
	for _, id := range ids {
		if topic, ok := topicMap[id]; ok {
			orderedTopics = append(orderedTopics, topic)
		}
	}

	return orderedTopics, nil
}
