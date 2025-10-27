/*
Package redis 负责 Redis 缓存连接管理
使用 go-redis/v9 库提供高性能的 Redis 客户端功能

主要功能：
1. 建立 Redis 连接池
2. 配置连接参数（地址、密码、数据库等）
3. 提供全局 Redis 客户端实例
4. 处理 Redis 连接的生命周期管理

技术栈：
- go-redis/v9: 高性能 Redis 客户端库
- Context: 支持上下文控制和超时管理

作者: DancingCircles
项目: Goweb-Frame
*/
package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9" // Redis 客户端库
	"github.com/spf13/viper"       // 配置管理库，用于读取 Redis 配置
	"go.uber.org/zap"              // 日志库，用于记录 Redis 操作日志
)

// rdb 全局 Redis 客户端实例
// 使用 *redis.Client 类型，提供完整的 Redis 操作功能
var rdb *redis.Client

// Init 初始化 Redis 连接
// 这个函数负责：
// 1. 从配置文件读取 Redis 连接参数
// 2. 创建 Redis 客户端实例
// 3. 测试连接是否有效
// 4. 配置连接池参数
//
// 返回值：
//   - error: 如果 Redis 连接失败，返回具体错误信息
func Init() (err error) {
	// ==================== 创建 Redis 客户端 ====================

	// 使用 redis.NewClient 创建客户端实例
	// Options 结构体包含所有连接配置参数
	rdb = redis.NewClient(&redis.Options{
		// 服务器地址，格式为 "host:port"
		Addr: fmt.Sprintf("%s:%d",
			viper.GetString("redis.host"), // Redis 服务器主机地址
			viper.GetInt("redis.port"),    // Redis 服务器端口
		),

		// Redis 服务器密码
		// 如果 Redis 没有设置密码，这个字段可以为空字符串
		Password: viper.GetString("redis.password"),

		// Redis 数据库编号 (0-15)
		// Redis 默认有 16 个数据库，编号从 0 到 15
		DB: viper.GetInt("redis.database"),

		// 连接池大小
		// 这个值决定了客户端可以同时维护的最大连接数
		PoolSize: viper.GetInt("redis.pool_size"),
	})

	// ==================== 测试连接有效性 ====================

	// 使用 Ping 命令测试连接是否正常
	// context.Background() 创建一个空的上下文
	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		zap.L().Error("Redis 连接失败", zap.Error(err))
		return err
	}

	// 记录成功连接的日志
	zap.L().Info("Redis 连接成功",
		zap.String("host", viper.GetString("redis.host")),
		zap.Int("port", viper.GetInt("redis.port")),
		zap.Int("database", viper.GetInt("redis.database")),
		zap.Int("pool_size", viper.GetInt("redis.pool_size")),
	)

	return nil
}

// Close 关闭 Redis 连接
// 这个函数会关闭 Redis 客户端和连接池中的所有连接
// 通常在应用程序退出时调用
func Close() {
	if rdb != nil {
		if err := rdb.Close(); err != nil {
			zap.L().Error("关闭 Redis 连接失败", zap.Error(err))
		} else {
			zap.L().Info("Redis 连接已关闭")
		}
	}
}

// GetRedis 获取 Redis 客户端实例
// 这个函数提供全局访问 Redis 客户端的方法
// 返回的是 *redis.Client 实例，可以用于执行 Redis 命令
//
// 使用示例：
//
//	rdb := redis.GetRedis()
//	err := rdb.Set(context.Background(), "key", "value", 0).Err()
//	val, err := rdb.Get(context.Background(), "key").Result()
//
// 返回值：
//   - *redis.Client: Redis 客户端实例
func GetRedis() *redis.Client {
	return rdb
}

// Redis 使用最佳实践说明：
//
// 1. 连接池配置：
//    - PoolSize: 根据应用并发量调整，通常设置为 CPU 核心数的 2-4 倍
//    - 过大：浪费内存资源
//    - 过小：可能成为性能瓶颈
//
// 2. 数据库选择：
//    - Redis 默认有 16 个数据库 (0-15)
//    - 不同功能使用不同数据库，便于管理和隔离
//    - 例如：0=缓存，1=会话，2=队列
//
// 3. 上下文使用：
//    - 所有 Redis 操作都需要传入 context
//    - 生产环境建议使用带超时的 context
//    - 例如：ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//
// 4. 错误处理：
//    - redis.Nil: 表示 key 不存在，这是正常情况
//    - 网络错误：需要重试机制
//    - 连接池耗尽：需要调整池大小或优化代码
//
// 5. 常用操作模式：
//    - 缓存：Set/Get 配合过期时间
//    - 会话：Hash 结构存储用户会话信息
//    - 计数器：Incr/Decr 原子操作
//    - 消息队列：List 的 Push/Pop 操作
