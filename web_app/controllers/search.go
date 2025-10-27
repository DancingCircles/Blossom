// Package controllers 处理HTTP请求的控制器
package controllers

import (
	"fmt"
	"net/http"
	"web_app/logic"
	"web_app/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SearchController 搜索控制器
type SearchController struct{}

// NewSearchController 创建搜索控制器
func NewSearchController() *SearchController {
	return &SearchController{}
}

// SearchTopics 搜索话题
// @Summary 搜索话题
// @Description 全文搜索话题，支持关键词、分类筛选和排序
// @Tags 搜索
// @Produce json
// @Param keyword query string false "搜索关键词"
// @Param category query string false "分类筛选"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param sort_by query string false "排序方式：created_at/view_count/comment_count" default(created_at)
// @Success 200 {object} models.Response{data=models.SearchResponse}
// @Router /api/v1/search [get]
func (sc *SearchController) SearchTopics(c *gin.Context) {
	// 1. 获取查询参数
	keyword := c.Query("keyword")
	category := c.Query("category")
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")
	sortBy := c.DefaultQuery("sort_by", "created_at")

	// 参数转换
	pageInt := 1
	pageSizeInt := 20
	if p, err := parseInt(page); err == nil {
		pageInt = p
	}
	if ps, err := parseInt(pageSize); err == nil {
		pageSizeInt = ps
	}

	// 2. 调用logic层搜索
	result, err := logic.SearchTopics(keyword, category, pageInt, pageSizeInt, sortBy)
	if err != nil {
		zap.L().Error("搜索话题失败",
			zap.Error(err),
			zap.String("keyword", keyword),
			zap.String("category", category))
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(models.CodeServerError, "搜索失败"))
		return
	}

	// 3. 返回搜索结果
	c.JSON(http.StatusOK, models.NewSuccessResponse(result))
}

// SuggestTopics 搜索建议
// @Summary 搜索建议
// @Description 根据前缀提供搜索建议
// @Tags 搜索
// @Produce json
// @Param prefix query string true "搜索前缀"
// @Success 200 {object} models.Response{data=[]string}
// @Router /api/v1/search/suggest [get]
func (sc *SearchController) SuggestTopics(c *gin.Context) {
	// 1. 获取前缀
	prefix := c.Query("prefix")
	if prefix == "" {
		c.JSON(http.StatusOK, models.NewSuccessResponse([]string{}))
		return
	}

	// 2. 调用logic层获取建议
	suggestions, err := logic.SuggestTopics(prefix)
	if err != nil {
		zap.L().Error("获取搜索建议失败", zap.Error(err), zap.String("prefix", prefix))
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(models.CodeServerError, "获取建议失败"))
		return
	}

	// 3. 返回建议列表
	c.JSON(http.StatusOK, models.NewSuccessResponse(suggestions))
}

// GetHotTopics 获取热门话题
// @Summary 获取分类热门话题
// @Description 获取指定分类的热门话题（按浏览量排序）
// @Tags 搜索
// @Produce json
// @Param category query string true "分类"
// @Param size query int false "数量" default(10)
// @Success 200 {object} models.Response{data=[]models.Topic}
// @Router /api/v1/search/hot [get]
func (sc *SearchController) GetHotTopics(c *gin.Context) {
	// 1. 获取参数
	category := c.Query("category")
	if category == "" {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.CodeInvalidParams, "分类不能为空"))
		return
	}

	size := 10
	if s := c.Query("size"); s != "" {
		if sInt, err := parseInt(s); err == nil {
			size = sInt
		}
	}

	// 2. 调用logic层
	topics, err := logic.GetHotTopicsByCategory(category, size)
	if err != nil {
		zap.L().Error("获取热门话题失败", zap.Error(err), zap.String("category", category))
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(models.CodeServerError, "获取热门话题失败"))
		return
	}

	// 3. 返回结果
	c.JSON(http.StatusOK, models.NewSuccessResponse(topics))
}

// GetCategoryStats 获取分类统计
// @Summary 获取分类统计
// @Description 统计各分类的话题数量
// @Tags 搜索
// @Produce json
// @Success 200 {object} models.Response{data=map[string]int64}
// @Router /api/v1/search/stats [get]
func (sc *SearchController) GetCategoryStats(c *gin.Context) {
	// 1. 调用logic层
	stats, err := logic.GetCategoryStats()
	if err != nil {
		zap.L().Error("获取分类统计失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(models.CodeServerError, "获取统计失败"))
		return
	}

	// 2. 返回统计结果
	c.JSON(http.StatusOK, models.NewSuccessResponse(stats))
}

// parseInt 字符串转整数的辅助函数
func parseInt(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}
