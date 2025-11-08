package utils

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var (
	ErrLockFailed   = errors.New("获取锁失败")
	ErrUnlockFailed = errors.New("释放锁失败")
)

// DistributedLock 分布式锁
type DistributedLock struct {
	client     *redis.Client
	key        string
	value      string
	expiration time.Duration
}

// NewDistributedLock 创建分布式锁实例
func NewDistributedLock(client *redis.Client, key string, expiration time.Duration) *DistributedLock {
	return &DistributedLock{
		client:     client,
		key:        key,
		value:      uuid.New().String(), // 使用UUID作为锁的值，防止误删
		expiration: expiration,
	}
}

// Lock 获取锁（使用 SET NX EX 原子操作）
func (l *DistributedLock) Lock(ctx context.Context) error {
	// SET key value NX EX seconds
	// NX: 只在键不存在时设置
	// EX: 设置过期时间（秒）
	success, err := l.client.SetNX(ctx, l.key, l.value, l.expiration).Result()
	if err != nil {
		return err
	}
	if !success {
		return ErrLockFailed
	}
	return nil
}

// Unlock 释放锁（使用Lua脚本保证原子性）
func (l *DistributedLock) Unlock(ctx context.Context) error {
	// Lua脚本：只有当锁的值等于自己的UUID时才删除
	// 防止删除别人的锁
	script := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`
	result, err := l.client.Eval(ctx, script, []string{l.key}, l.value).Result()
	if err != nil {
		return err
	}
	if result == int64(0) {
		return ErrUnlockFailed
	}
	return nil
}

// TryLock 尝试获取锁（带重试）
func (l *DistributedLock) TryLock(ctx context.Context, retryCount int, retryDelay time.Duration) error {
	for i := 0; i < retryCount; i++ {
		err := l.Lock(ctx)
		if err == nil {
			return nil
		}
		if err != ErrLockFailed {
			return err
		}
		// 等待后重试
		time.Sleep(retryDelay)
	}
	return ErrLockFailed
}

// WithLock 使用锁执行函数（自动加锁/解锁）
func WithLock(ctx context.Context, client *redis.Client, key string, expiration time.Duration, fn func() error) error {
	lock := NewDistributedLock(client, key, expiration)

	// 获取锁
	if err := lock.Lock(ctx); err != nil {
		return err
	}

	// 确保函数返回后释放锁
	defer func() {
		if err := lock.Unlock(ctx); err != nil {
			// 记录日志，但不影响主流程
		}
	}()

	// 执行业务逻辑
	return fn()
}
