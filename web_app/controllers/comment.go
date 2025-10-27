// Package controllers 处理HTTP请求的控制器
package controllers

import (
	"net/http"
	"strconv"
	"web_app/logic"
	"web_app/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CommentController 评论控制器
type CommentController struct{}

// NewCommentController 创建评论控制器
func NewCommentController() *CommentController {
	return &CommentController{}
}

// CreateComment 创建评论
// @Summary 创建评论
// @Description 对话题发表评论
// @Tags 评论
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "话题ID"
// @Param request body models.CreateCommentRequest true "评论信息"
// @Success 200 {object} models.Response
// @Router /api/v1/topics/{id}/comments [post]
func (cc *CommentController) CreateComment(c *gin.Context) {
	// 1. 获取话题ID
	topicID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.CodeInvalidParams, "无效的话题ID"))
		return
	}

	// 2. 绑定并验证请求参数
	var req models.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.L().Error("参数验证失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.CodeInvalidParams, "参数错误: "+err.Error()))
		return
	}

	// 3. 从context获取当前用户ID
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.CodeUnauthorized, "用户未登录"))
		return
	}
	userID, ok := userIDVal.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(models.CodeServerError, "用户ID格式解析错误"))
		return
	}

	// 4. 调用逻辑层创建评论
	if err := logic.CreateComment(userID, topicID, &req); err != nil {
		zap.L().Error("创建评论失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(models.CodeServerError, err.Error()))
		return
	}

	// 5. 返回成功响应
	c.JSON(http.StatusOK, models.NewSuccessResponse(gin.H{
		"message": "评论发表成功",
	}))
}

// GetComments 获取评论列表
// @Summary 获取话题评论列表
// @Description 分页获取话题的评论列表
// @Tags 评论
// @Produce json
// @Param id path int true "话题ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} models.Response{data=models.CommentListResponse}
// @Router /api/v1/topics/{id}/comments [get]
func (cc *CommentController) GetComments(c *gin.Context) {
	// 1. 获取话题ID
	topicID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.CodeInvalidParams, "无效的话题ID"))
		return
	}

	// 2. 绑定查询参数
	var req models.GetCommentsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.CodeInvalidParams, "参数错误"))
		return
	}

	// 设置默认值
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	// 3. 通过逻辑层查询数据
	commentList, total, err := logic.GetCommentsByTopicID(topicID, &req)
	if err != nil {
		// 区分"话题不存在"和其他错误
		if err.Error() == "话题不存在" {
			zap.L().Warn("话题不存在", zap.Int64("topic_id", topicID))
			c.JSON(http.StatusNotFound, models.NewErrorResponse(models.CodeNotFound, "话题不存在"))
			return
		}
		zap.L().Error("获取评论列表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(models.CodeServerError, "获取评论列表失败"))
		return
	}

	// 4. 处理空列表情况
	if commentList == nil {
		commentList = []*models.Comment{}
	}

	// 5. 计算分页信息
	totalPages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
	hasMore := req.Page < totalPages

	// 6. 返回响应
	c.JSON(http.StatusOK, models.NewSuccessResponse(models.CommentListResponse{
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
		HasMore:    hasMore,
		Comments:   commentList,
	}))
}

// DeleteComment 删除评论
// @Summary 删除评论
// @Description 删除自己的评论
// @Tags 评论
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "评论ID"
// @Success 200 {object} models.Response
// @Router /api/v1/comments/{id} [delete]
func (cc *CommentController) DeleteComment(c *gin.Context) {
	// 1. 获取评论ID
	commentID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.CodeInvalidParams, "无效的评论ID"))
		return
	}

	// 2. 从context获取当前用户ID
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.CodeUnauthorized, "用户未登录"))
		return
	}
	userID, ok := userIDVal.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(models.CodeServerError, "用户ID格式解析错误"))
		return
	}

	// 3. 调用逻辑层删除评论
	if err := logic.DeleteComment(userID, commentID); err != nil {
		zap.L().Error("删除评论失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(models.CodeServerError, err.Error()))
		return
	}

	// 4. 返回成功响应
	c.JSON(http.StatusOK, models.NewSuccessResponse(gin.H{
		"message": "评论删除成功",
	}))
}
