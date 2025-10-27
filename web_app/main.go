// Package main 是Go Web应用程序入口点
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"web_app/dao/elasticsearch"
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/logger"
	"web_app/routes"
	"web_app/settings"
	"web_app/utils"

	_ "net/http/pprof" // 性能监控

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// @title Bullbell Forum API
// @version 1.0
// @description 基于Go+Gin的论坛系统API，支持用户注册登录、话题发布、评论互动等功能
// @contact.name DancingCircles
// @contact.url https://github.com/DancingCircles/Bullbell
// @host localhost:8082
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description JWT token，格式：Bearer {token}

// main 应用程序入口点
func main() {
	// 初始化配置
	if err := settings.Init(); err != nil {
		fmt.Printf("初始化配置失败, 错误:%v\n", err)
		return
	}

	// 初始化日志
	if err := logger.Init(); err != nil {
		fmt.Printf("初始化日志系统失败, 错误:%v\n", err)
		return
	}
	defer func() {
		if err := zap.L().Sync(); err != nil {
			// 忽略 sync 错误（在某些系统上会出现 invalid argument 错误）
			_ = err
		}
	}()
	zap.L().Debug("日志系统初始化成功...")

	// 初始化雪花算法ID生成器
	if err := utils.InitSnowflake(); err != nil {
		fmt.Printf("初始化雪花算法失败, 错误:%v\n", err)
		return
	}
	zap.L().Debug("雪花算法ID生成器初始化成功...")

	// 初始化MySQL
	if err := mysql.Init(); err != nil {
		fmt.Printf("初始化MySQL失败, 错误:%v\n", err)
		return
	}
	defer mysql.Close()

	// 初始化Redis
	if err := redis.Init(); err != nil {
		fmt.Printf("初始化Redis失败, 错误:%v\n", err)
		return
	}
	defer redis.Close()

	// 初始化Elasticsearch
	if err := elasticsearch.Init(); err != nil {
		fmt.Printf("初始化Elasticsearch失败, 错误:%v\n", err)
		return
	}
	defer elasticsearch.Close()
	zap.L().Debug("Elasticsearch初始化成功...")

	// 启动pprof性能监控服务（仅开发环境）
	if viper.GetString("app.mode") == "dev" {
		go func() {
			zap.L().Info("正在启动pprof性能监控服务", zap.String("addr", "http://localhost:6060/debug/pprof"))
			if err := http.ListenAndServe("localhost:6060", nil); err != nil {
				zap.L().Error("pprof服务启动失败", zap.Error(err))
			}
		}()
	}

	// 初始化路由
	r := routes.SetupRouter()

	// 配置HTTP服务器
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("app.port")),
		Handler: r,
	}

	// 启动HTTP服务器
	go func() {
		zap.L().Info("正在启动HTTP服务器", zap.Int("port", viper.GetInt("app.port")))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务器启动失败: %s\n", err)
		}
	}()

	// 优雅关机
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	zap.L().Info("收到关闭信号，开始优雅关机...")

	// 5秒超时关机
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 关闭服务器
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("服务器强制关闭", zap.Error(err))
	}

	zap.L().Info("服务器已安全退出")
}
