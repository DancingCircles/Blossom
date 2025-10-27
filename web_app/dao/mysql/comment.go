package mysql

import (
	"database/sql"
	"errors"
	"web_app/models"
)

// InsertComment 插入评论
func InsertComment(comment *models.Comment) error {
	sqlStr := "INSERT INTO comments (id, topic_id, user_id, content, parent_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)"
	_, err := db.Exec(sqlStr, comment.ID, comment.TopicID, comment.UserID, comment.Content, comment.ParentID, comment.CreatedAt, comment.UpdatedAt)
	return err
}

// GetCommentsByTopicID 根据话题ID获取评论列表（分页）
func GetCommentsByTopicID(topicID int64, req *models.GetCommentsRequest) ([]*models.Comment, int64, error) {
	// 查询总数
	countSQL := "SELECT COUNT(*) FROM comments WHERE topic_id = ?"
	var total int64
	if err := db.Get(&total, countSQL, topicID); err != nil {
		return nil, 0, err
	}

	// 查询评论列表（JOIN users 表获取用户名）
	offset := (req.Page - 1) * req.PageSize
	listSQL := `
		SELECT c.*, u.username 
		FROM comments c
		LEFT JOIN users u ON c.user_id = u.id
		WHERE c.topic_id = ?
		ORDER BY c.created_at ASC
		LIMIT ? OFFSET ?
	`

	var commentList []*models.Comment
	if err := db.Select(&commentList, listSQL, topicID, req.PageSize, offset); err != nil {
		return nil, 0, err
	}

	return commentList, total, nil
}

// GetCommentByID 根据ID获取评论详情
func GetCommentByID(commentID int64) (*models.Comment, error) {
	sqlStr := `
		SELECT c.*, u.username 
		FROM comments c
		LEFT JOIN users u ON c.user_id = u.id
		WHERE c.id = ?
	`
	var comment models.Comment
	if err := db.Get(&comment, sqlStr, commentID); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("评论不存在")
		}
		return nil, err
	}
	return &comment, nil
}

// DeleteComment 删除评论
func DeleteComment(commentID int64) error {
	sqlStr := "DELETE FROM comments WHERE id = ?"
	result, err := db.Exec(sqlStr, commentID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("评论不存在")
	}

	return nil
}

// UpdateTopicCommentCount 更新话题评论数
func UpdateTopicCommentCount(topicID int64, delta int) error {
	sqlStr := "UPDATE topics SET comment_count = comment_count + ? WHERE id = ?"
	_, err := db.Exec(sqlStr, delta, topicID)
	return err
}

// CountCommentsByTopicID 统计话题的评论数
func CountCommentsByTopicID(topicID int64) (int64, error) {
	var count int64
	sqlStr := "SELECT COUNT(*) FROM comments WHERE topic_id = ?"
	if err := db.Get(&count, sqlStr, topicID); err != nil {
		return 0, err
	}
	return count, nil
}
