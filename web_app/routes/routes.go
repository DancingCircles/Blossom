// Package routes 负责HTTP路由配置和中间件管理
package routes

import (
	"net/http"
	"web_app/controllers"
	_ "web_app/docs" // Swagger文档
	"web_app/logger"
	"web_app/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter 配置和初始化路由
func SetupRouter() *gin.Engine {
	// 创建路由引擎
	r := gin.New()

	// ========== 全局中间件 ==========
	r.Use(
		logger.GinLogger(),               // 日志中间件
		logger.GinRecovery(true),         // 恢复中间件（panic恢复）
		middleware.CORS(),                // 跨域中间件
		middleware.IPRateLimit(100, 200), // IP限流：每秒100个请求，桶容量200
	)

	// ========== 健康检查接口 ==========
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
			"status":  "healthy",
			"service": "Blossom Forum",
		})
	})

	// ========== Swagger API文档 ==========
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ========== 初始化控制器 ==========
	userCtrl := controllers.NewUserController()
	topicCtrl := controllers.NewTopicController()
	commentCtrl := controllers.NewCommentController()
	searchCtrl := controllers.NewSearchController()
	adminCtrl := controllers.NewAdminController()

	// ========== API 路由组 ==========
	api := r.Group("/api")
	{
		// v1 版本API
		v1 := api.Group("/v1")
		{
			// ===== 公开接口（无需登录） =====
			v1.POST("/register", userCtrl.Register) // 用户注册
			v1.POST("/login", userCtrl.Login)       // 用户登录

			// 话题列表（无需登录也可以查看）
			v1.GET("/topics", topicCtrl.GetTopics)                  // 获取话题列表
			v1.GET("/topics/:id", topicCtrl.GetTopicByID)           // 获取话题详情
			v1.GET("/topics/:id/comments", commentCtrl.GetComments) // 获取话题评论列表

			// 搜索相关（无需登录）
			v1.GET("/search", searchCtrl.SearchTopics)           // 搜索话题
			v1.GET("/search/suggest", searchCtrl.SuggestTopics)  // 搜索建议
			v1.GET("/search/hot", searchCtrl.GetHotTopics)       // 热门话题
			v1.GET("/search/stats", searchCtrl.GetCategoryStats) // 分类统计

			// 管理接口（临时公开，生产环境应加权限验证）
			v1.POST("/admin/sync-es", adminCtrl.SyncToES) // 同步数据到ES

			// ===== 需要登录的接口 =====
			// 使用JWT中间件保护
			auth := v1.Group("")
			auth.Use(middleware.JWTAuth())
			{
				// 用户相关
				auth.GET("/user/info", userCtrl.GetUserInfo) // 获取当前用户信息

				// 话题相关
				auth.POST("/topics", topicCtrl.CreateTopic)        // 创建话题
				auth.POST("/topics/:id/vote", topicCtrl.VoteTopic) // 给话题投票

				// 评论相关
				auth.POST("/topics/:id/comments", commentCtrl.CreateComment) // 发表评论
				auth.DELETE("/comments/:id", commentCtrl.DeleteComment)      // 删除评论
			}
		}
	}

	// 返回配置完成的路由引擎
	return r
}
