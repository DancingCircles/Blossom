package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"web_app/models"

	"github.com/jmoiron/sqlx"
)

// InsertTopic 插入话题
func InsertTopic(topic *models.Topic) error {
	sqlStr := "INSERT INTO topics (id, user_id, title, content, category, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)"
	_, err := db.Exec(sqlStr, topic.ID, topic.UserID, topic.Title, topic.Content, topic.Category, topic.CreatedAt, topic.UpdatedAt)
	return err
}

// GetTopics 获取话题列表（带分页、排序、筛选）
func GetTopics(req *models.GetTopicsRequest) ([]*models.Topic, int64, error) {
	// 构建 WHERE 条件
	whereClause := "WHERE 1=1"
	args := []interface{}{}

	if req.Category != "" {
		whereClause += " AND t.category = ?"
		args = append(args, req.Category)
	}

	// 构建 ORDER BY 子句
	var orderBy string
	switch req.Sort {
	case "new":
		orderBy = "ORDER BY t.created_at DESC"
	case "like":
		orderBy = "ORDER BY t.like_count DESC"
	case "hot":
		fallthrough
	default:
		// 热度算法：综合点赞数、评论数、浏览数
		orderBy = "ORDER BY (t.like_count * 3 + t.comment_count * 2 + t.view_count) DESC"
	}

	// 查询总数
	countSQL := fmt.Sprintf("SELECT COUNT(*) FROM topics t %s", whereClause)
	var total int64
	if err := db.Get(&total, countSQL, args...); err != nil {
		return nil, 0, err
	}

	// 查询话题列表（JOIN users 表获取用户名）
	offset := (req.Page - 1) * req.PageSize
	listSQL := fmt.Sprintf(`
		SELECT t.*, u.username 
		FROM topics t
		LEFT JOIN users u ON t.user_id = u.id
		%s
		%s
		LIMIT ? OFFSET ?
	`, whereClause, orderBy)

	args = append(args, req.PageSize, offset)

	var topicList []*models.Topic
	if err := db.Select(&topicList, listSQL, args...); err != nil {
		return nil, 0, err
	}

	return topicList, total, nil
}

// GetTopicByID 根据ID获取话题详情
func GetTopicByID(topicID int64) (*models.Topic, error) {
	sqlStr := `
		SELECT t.*, u.username 
		FROM topics t
		LEFT JOIN users u ON t.user_id = u.id
		WHERE t.id = ?
	`
	var topic models.Topic
	if err := db.Get(&topic, sqlStr, topicID); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("话题不存在")
		}
		return nil, err
	}
	return &topic, nil
}

// GetAllTopics 获取所有话题（用于同步到ES）
func GetAllTopics() ([]*models.Topic, error) {
	sqlStr := `
		SELECT t.*, u.username 
		FROM topics t
		LEFT JOIN users u ON t.user_id = u.id
		ORDER BY t.created_at DESC
	`
	var topics []*models.Topic
	if err := db.Select(&topics, sqlStr); err != nil {
		return nil, err
	}
	return topics, nil
}

// UpdateTopicViewCount 增加话题浏览数
func UpdateTopicViewCount(topicID int64) error {
	sqlStr := "UPDATE topics SET view_count = view_count + 1 WHERE id = ?"
	_, err := db.Exec(sqlStr, topicID)
	return err
}

// GetUserVote 获取用户对话题的投票状态
func GetUserVote(userID, topicID int64) (*models.Vote, error) {
	sqlStr := "SELECT * FROM votes WHERE user_id = ? AND topic_id = ?"
	var vote models.Vote
	if err := db.Get(&vote, sqlStr, userID, topicID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // 未投票
		}
		return nil, err
	}
	return &vote, nil
}

// InsertVote 插入投票记录
func InsertVote(vote *models.Vote) error {
	sqlStr := "INSERT INTO votes (id, user_id, topic_id, vote_type) VALUES (?, ?, ?, ?)"
	_, err := db.Exec(sqlStr, vote.ID, vote.UserID, vote.TopicID, vote.VoteType)
	return err
}

// UpdateVote 更新投票记录
func UpdateVote(userID, topicID int64, voteType int) error {
	sqlStr := "UPDATE votes SET vote_type = ? WHERE user_id = ? AND topic_id = ?"
	_, err := db.Exec(sqlStr, voteType, userID, topicID)
	return err
}

// DeleteVote 删除投票记录
func DeleteVote(userID, topicID int64) error {
	sqlStr := "DELETE FROM votes WHERE user_id = ? AND topic_id = ?"
	_, err := db.Exec(sqlStr, userID, topicID)
	return err
}

// UpdateTopicLikeCount 更新话题点赞数
func UpdateTopicLikeCount(topicID int64, delta int) error {
	sqlStr := "UPDATE topics SET like_count = like_count + ? WHERE id = ?"
	_, err := db.Exec(sqlStr, delta, topicID)
	return err
}

// UpdateTopicDislikeCount 更新话题点踩数
func UpdateTopicDislikeCount(topicID int64, delta int) error {
	sqlStr := "UPDATE topics SET dislike_count = dislike_count + ? WHERE id = ?"
	_, err := db.Exec(sqlStr, delta, topicID)
	return err
}

// GetRecentTopics 获取最近的话题列表（用于热度排名计算）
func GetRecentTopics(limit int) ([]*models.Topic, error) {
	sqlStr := `SELECT t.id, t.user_id, u.username, t.title, t.content, t.category, 
	                  t.like_count, t.dislike_count, t.comment_count, t.view_count,
	                  t.created_at, t.updated_at
	           FROM topics t
	           LEFT JOIN users u ON t.user_id = u.id
	           ORDER BY t.created_at DESC
	           LIMIT ?`

	var topics []*models.Topic
	err := db.Select(&topics, sqlStr, limit)
	if err != nil {
		return nil, err
	}

	return topics, nil
}

// GetTopicsByIDs 根据ID列表批量获取话题
func GetTopicsByIDs(ids []int64) ([]*models.Topic, error) {
	if len(ids) == 0 {
		return []*models.Topic{}, nil
	}

	// 构建IN查询
	query, args, err := sqlx.In(`SELECT t.id, t.user_id, u.username, t.title, t.content, t.category,
	                                    t.like_count, t.dislike_count, t.comment_count, t.view_count,
	                                    t.created_at, t.updated_at
	                             FROM topics t
	                             LEFT JOIN users u ON t.user_id = u.id
	                             WHERE t.id IN (?)`, ids)
	if err != nil {
		return nil, err
	}

	query = db.Rebind(query)
	var topics []*models.Topic
	err = db.Select(&topics, query, args...)
	if err != nil {
		return nil, err
	}

	return topics, nil
}
