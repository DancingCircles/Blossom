package logic

import (
	"errors"
	"time"
	"web_app/dao/elasticsearch"
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/models"
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

	// 异步索引到Elasticsearch
	go func() {
		if err := elasticsearch.IndexTopic(topic); err != nil {
			zap.L().Warn("索引话题到ES失败", zap.Error(err), zap.Int64("topic_id", topicID))
		}
	}()

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
