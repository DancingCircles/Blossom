// Package controllers 管理员控制器
package controllers

import (
	"net/http"
	"web_app/dao/elasticsearch"
	"web_app/dao/mysql"
	"web_app/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AdminController 管理员控制器
type AdminController struct{}

// NewAdminController 创建管理员控制器
func NewAdminController() *AdminController {
	return &AdminController{}
}

// SyncToES 同步数据到Elasticsearch
// @Summary 同步数据到ES
// @Description 将MySQL中的所有话题同步到Elasticsearch
// @Tags 管理
// @Produce json
// @Success 200 {object} models.Response
// @Router /api/v1/admin/sync-es [post]
func (ac *AdminController) SyncToES(c *gin.Context) {
	zap.L().Info("开始同步数据到Elasticsearch")

	// 1. 从MySQL获取所有话题
	topics, err := mysql.GetAllTopics()
	if err != nil {
		zap.L().Error("获取话题失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(models.CodeServerError, "获取话题失败"))
		return
	}

	if len(topics) == 0 {
		c.JSON(http.StatusOK, models.NewSuccessResponse(gin.H{
			"message": "没有话题需要同步",
			"count":   0,
		}))
		return
	}

	// 2. 批量索引到ES
	err = elasticsearch.BulkIndexTopics(topics)
	if err != nil {
		zap.L().Error("批量索引失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(models.CodeServerError, "同步失败"))
		return
	}

	zap.L().Info("同步完成", zap.Int("count", len(topics)))

	c.JSON(http.StatusOK, models.NewSuccessResponse(gin.H{
		"message": "同步成功",
		"count":   len(topics),
	}))
}
