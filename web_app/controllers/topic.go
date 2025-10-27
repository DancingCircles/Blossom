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

// TopicController 话题控制器
type TopicController struct{}

// NewTopicController 创建话题控制器
func NewTopicController() *TopicController {
	return &TopicController{}
}

// CreateTopic 创建话题
// @Summary 创建话题
// @Description 发布新话题
// @Tags 话题
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body models.CreateTopicRequest true "话题信息"
// @Success 200 {object} models.Response
// @Router /api/v1/topics [post]
func (tc *TopicController) CreateTopic(c *gin.Context) {
	// 1. 绑定并验证请求参数
	var req models.CreateTopicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.L().Error("参数验证失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.CodeInvalidParams, "参数错误: "+err.Error()))
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

	//3. 调用逻辑层插入话题
	if err := logic.CreateTopic(userID, &req); err != nil {
		zap.L().Error("创建话题失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(models.CodeServerError, err.Error()))
		return
	}

	// 4. 返回成功响应
	c.JSON(http.StatusOK, models.NewSuccessResponse(gin.H{
		"message": "话题创建成功",
	}))
}

// GetTopics 获取话题列表
// @Summary 获取话题列表
// @Description 分页获取话题列表，支持排序和筛选
// @Tags 话题
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param sort query string false "排序方式：hot/new/like" default(hot)
// @Param category query string false "分类筛选"
// @Success 200 {object} models.Response{data=models.TopicListResponse}
// @Router /api/v1/topics [get]
func (tc *TopicController) GetTopics(c *gin.Context) {
	// 1. 绑定查询参数
	var req models.GetTopicsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.CodeInvalidParams, "参数错误"))
		return
	}

	// 设置默认值
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 10
	}
	if req.Sort == "" {
		req.Sort = "hot"
	}

	// 2. 通过逻辑层查询数据
	topicList, total, err := logic.GetTopics(&req)
	if err != nil {
		zap.L().Error("获取话题列表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(models.CodeServerError, "获取话题列表失败"))
		return
	}

	// 3. 处理空列表情况
	if topicList == nil {
		topicList = []*models.Topic{}
	}

	// 4. 计算分页信息
	totalPages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
	hasMore := req.Page < totalPages

	// 5. 返回响应
	c.JSON(http.StatusOK, models.NewSuccessResponse(models.TopicListResponse{
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
		HasMore:    hasMore,
		Topics:     topicList,
	}))
}

// GetTopicByID 获取话题详情
// @Summary 获取话题详情
// @Description 根据ID获取话题详细信息
// @Tags 话题
// @Produce json
// @Param id path int true "话题ID"
// @Success 200 {object} models.Response{data=models.Topic}
// @Router /api/v1/topics/{id} [get]
func (tc *TopicController) GetTopicByID(c *gin.Context) {
	// 1. 获取话题ID
	topicID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.CodeInvalidParams, "无效的话题ID"))
		return
	}

	// 2. 调用逻辑层查询话题详情
	topic, err := logic.GetTopicByID(topicID)
	if err != nil {
		// 区分"话题不存在"和其他错误
		if err.Error() == "话题不存在" {
			zap.L().Warn("话题不存在", zap.Int64("topic_id", topicID))
			c.JSON(http.StatusNotFound, models.NewErrorResponse(models.CodeNotFound, "话题不存在"))
			return
		}
		zap.L().Error("获取话题详情失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(models.CodeServerError, "获取话题详情失败"))
		return
	}

	// 3. 返回话题详情
	c.JSON(http.StatusOK, models.NewSuccessResponse(topic))
}

// VoteTopic 给话题投票（点赞/点踩）
// @Summary 给话题投票
// @Description 对话题进行点赞或点踩
// @Tags 话题
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "话题ID"
// @Param type query string true "投票类型：like/dislike"
// @Success 200 {object} models.Response
// @Router /api/v1/topics/{id}/vote [post]
func (tc *TopicController) VoteTopic(c *gin.Context) {
	// 1. 获取话题ID
	topicID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.CodeInvalidParams, "无效的话题ID"))
		return
	}

	// 2. 获取投票类型
	voteType := c.Query("type") // "like" 或 "dislike"
	if voteType != "like" && voteType != "dislike" {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.CodeInvalidParams, "无效的投票类型"))
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

	// 4. 调用逻辑层进行投票处理
	if err := logic.VoteTopic(userID, topicID, voteType); err != nil {
		zap.L().Error("话题投票失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(models.CodeServerError, err.Error()))
		return
	}

	// 5. 返回成功响应
	c.JSON(http.StatusOK, models.NewSuccessResponse(gin.H{
		"message": "投票成功",
		"type":    voteType,
	}))
}
