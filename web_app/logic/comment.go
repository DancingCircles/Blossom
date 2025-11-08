package logic

import (
	"context"
	"errors"
	"fmt"
	"time"
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/models"
	"web_app/utils"

	"go.uber.org/zap"
)

// CreateComment 创建评论
func CreateComment(userID, topicID int64, req *models.CreateCommentRequest) error {
	// 使用分布式锁防止短时间内重复提交评论
	ctx := context.Background()
	lockKey := fmt.Sprintf("lock:comment:%d:%d", topicID, userID)

	// 使用 WithLock 自动管理锁的获取和释放（2秒超时，防止用户短时间内重复提交）
	return utils.WithLock(ctx, redis.GetClient(), lockKey, 2*time.Second, func() error {
		// 1. 验证话题是否存在
		_, err := mysql.GetTopicByID(topicID)
		if err != nil {
			return errors.New("话题不存在")
		}

		// 2. 如果有父评论ID，验证父评论是否存在且属于同一话题
		if req.ParentID != nil {
			parentComment, err := mysql.GetCommentByID(*req.ParentID)
			if err != nil {
				return errors.New("父评论不存在")
			}
			if parentComment.TopicID != topicID {
				return errors.New("父评论不属于当前话题")
			}
		}

		// 3. 生成雪花算法ID
		commentID := utils.GenerateID()

		// 4. 插入评论（显式设置创建时间为当前时间）
		now := time.Now()
		err = mysql.InsertComment(&models.Comment{
			ID:        commentID,
			TopicID:   topicID,
			UserID:    userID,
			Content:   req.Content,
			ParentID:  req.ParentID,
			CreatedAt: now,
			UpdatedAt: now,
		})
		if err != nil {
			zap.L().Error("插入评论失败", zap.Error(err))
			return errors.New("插入评论失败")
		}

		// 5. 更新话题评论数（异步处理，不影响主流程）
		go func() {
			if err := mysql.UpdateTopicCommentCount(topicID, 1); err != nil {
				zap.L().Warn("更新话题评论数失败", zap.Error(err))
			}
		}()

		return nil
	})
}

// GetCommentsByTopicID 获取话题评论列表
func GetCommentsByTopicID(topicID int64, req *models.GetCommentsRequest) ([]*models.Comment, int64, error) {
	// 1. 验证话题是否存在
	_, err := mysql.GetTopicByID(topicID)
	if err != nil {
		return nil, 0, errors.New("话题不存在")
	}

	// 2. 调用dao层查询评论
	commentList, total, err := mysql.GetCommentsByTopicID(topicID, req)
	if err != nil {
		zap.L().Error("查询评论失败", zap.Error(err))
		return nil, 0, errors.New("查询评论失败")
	}

	return commentList, total, nil
}

// DeleteComment 删除评论（仅允许作者删除）
func DeleteComment(userID, commentID int64) error {
	// 1. 查询评论是否存在
	comment, err := mysql.GetCommentByID(commentID)
	if err != nil {
		return errors.New("评论不存在")
	}

	// 2. 验证是否是评论作者
	if comment.UserID != userID {
		return errors.New("无权删除该评论")
	}

	// 3. 删除评论
	if err := mysql.DeleteComment(commentID); err != nil {
		zap.L().Error("删除评论失败", zap.Error(err))
		return errors.New("删除评论失败")
	}

	// 4. 更新话题评论数（异步处理，不影响主流程）
	go func() {
		if err := mysql.UpdateTopicCommentCount(comment.TopicID, -1); err != nil {
			zap.L().Warn("更新话题评论数失败", zap.Error(err))
		}
	}()

	return nil
}
