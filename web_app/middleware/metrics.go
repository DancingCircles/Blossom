// Package middleware 提供HTTP中间件
package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP请求总数
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "HTTP请求总数",
		},
		[]string{"method", "path", "status"},
	)

	// HTTP请求延迟直方图
	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP请求延迟（秒）",
			Buckets: prometheus.DefBuckets, // 默认桶：0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10
		},
		[]string{"method", "path", "status"},
	)

	// 正在处理的HTTP请求数
	httpRequestsInProgress = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_progress",
			Help: "当前正在处理的HTTP请求数",
		},
	)
)

// PrometheusMetrics Prometheus指标采集中间件
func PrometheusMetrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过metrics端点本身，避免自我监控
		if c.Request.URL.Path == "/metrics" {
			c.Next()
			return
		}

		// 请求开始时增加计数
		httpRequestsInProgress.Inc()
		start := time.Now()

		// 处理请求
		c.Next()

		// 请求结束时减少计数
		httpRequestsInProgress.Dec()

		// 记录请求信息
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())
		method := c.Request.Method
		path := c.FullPath()

		// 如果path为空（404等情况），使用原始路径
		if path == "" {
			path = c.Request.URL.Path
		}

		// 更新指标
		httpRequestsTotal.WithLabelValues(method, path, status).Inc()
		httpRequestDuration.WithLabelValues(method, path, status).Observe(duration)
	}
}


