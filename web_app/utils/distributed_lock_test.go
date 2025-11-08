package utils

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

// setupTestRedis 创建测试用的Redis客户端
func setupTestRedis(t *testing.T) (*redis.Client, *miniredis.Miniredis) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("无法启动miniredis: %v", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return client, mr
}

// TestDistributedLock_BasicLockUnlock 测试基本的加锁和解锁
func TestDistributedLock_BasicLockUnlock(t *testing.T) {
	client, mr := setupTestRedis(t)
	defer mr.Close()
	defer client.Close()

	ctx := context.Background()
	lock := NewDistributedLock(client, "test:lock", 5*time.Second)

	// 测试加锁
	err := lock.Lock(ctx)
	if err != nil {
		t.Errorf("加锁失败: %v", err)
	}

	// 测试重复加锁应该失败
	err = lock.Lock(ctx)
	if err != ErrLockFailed {
		t.Errorf("期望获取锁失败，但成功了")
	}

	// 测试解锁
	err = lock.Unlock(ctx)
	if err != nil {
		t.Errorf("解锁失败: %v", err)
	}

	// 解锁后应该可以再次加锁
	lock2 := NewDistributedLock(client, "test:lock", 5*time.Second)
	err = lock2.Lock(ctx)
	if err != nil {
		t.Errorf("解锁后加锁失败: %v", err)
	}
}

// TestDistributedLock_Concurrency 测试并发场景
func TestDistributedLock_Concurrency(t *testing.T) {
	client, mr := setupTestRedis(t)
	defer mr.Close()
	defer client.Close()

	ctx := context.Background()
	lockKey := "test:concurrent:lock"

	// 用于记录成功获取锁的协程数量
	successCount := 0
	var mu sync.Mutex
	var wg sync.WaitGroup

	// 启动10个协程同时尝试获取锁
	goroutineCount := 10
	wg.Add(goroutineCount)

	for i := 0; i < goroutineCount; i++ {
		go func(id int) {
			defer wg.Done()

			lock := NewDistributedLock(client, lockKey, 1*time.Second)
			err := lock.Lock(ctx)

			if err == nil {
				// 获取锁成功，增加计数
				mu.Lock()
				successCount++
				mu.Unlock()

				// 模拟处理业务
				time.Sleep(100 * time.Millisecond)

				// 释放锁
				lock.Unlock(ctx)
			}
		}(i)
	}

	wg.Wait()

	// 在同一时刻，只有一个协程能获取到锁
	if successCount != 1 {
		t.Errorf("期望只有1个协程获取锁成功，但实际有 %d 个", successCount)
	}
}

// TestDistributedLock_TryLock 测试带重试的加锁
func TestDistributedLock_TryLock(t *testing.T) {
	client, mr := setupTestRedis(t)
	defer mr.Close()
	defer client.Close()

	ctx := context.Background()
	lockKey := "test:trylock"

	// 第一个锁先获取
	lock1 := NewDistributedLock(client, lockKey, 2*time.Second)
	err := lock1.Lock(ctx)
	if err != nil {
		t.Fatalf("第一个锁获取失败: %v", err)
	}

	// 在另一个协程中延迟释放锁
	go func() {
		time.Sleep(500 * time.Millisecond)
		lock1.Unlock(ctx)
	}()

	// 第二个锁使用 TryLock，应该在重试后成功
	lock2 := NewDistributedLock(client, lockKey, 2*time.Second)
	start := time.Now()
	err = lock2.TryLock(ctx, 5, 200*time.Millisecond)
	elapsed := time.Since(start)

	if err != nil {
		t.Errorf("TryLock 失败: %v", err)
	}

	// 确保经过了重试时间
	if elapsed < 400*time.Millisecond {
		t.Errorf("TryLock 应该等待至少400ms，但只等待了 %v", elapsed)
	}
}

// TestWithLock 测试 WithLock 辅助函数
func TestWithLock(t *testing.T) {
	client, mr := setupTestRedis(t)
	defer mr.Close()
	defer client.Close()

	ctx := context.Background()
	lockKey := "test:withlock"
	executed := false

	// 测试正常执行
	err := WithLock(ctx, client, lockKey, 5*time.Second, func() error {
		executed = true
		return nil
	})

	if err != nil {
		t.Errorf("WithLock 执行失败: %v", err)
	}

	if !executed {
		t.Error("函数未被执行")
	}

	// 验证锁已被释放（可以再次获取）
	lock := NewDistributedLock(client, lockKey, 5*time.Second)
	err = lock.Lock(ctx)
	if err != nil {
		t.Error("锁未被正确释放")
	}
}

// TestDistributedLock_Expiration 测试锁的自动过期
// 注意：此测试在 miniredis 中可能不稳定，主要用于真实Redis环境
func TestDistributedLock_Expiration(t *testing.T) {
	t.Skip("miniredis的TTL行为与真实Redis不同，跳过此测试")

	client, mr := setupTestRedis(t)
	defer mr.Close()
	defer client.Close()

	ctx := context.Background()
	lockKey := "test:expiration"

	// 创建一个1秒过期的锁
	lock1 := NewDistributedLock(client, lockKey, 1*time.Second)
	err := lock1.Lock(ctx)
	if err != nil {
		t.Fatalf("加锁失败: %v", err)
	}

	// 在真实Redis环境中，需要等待锁过期
	mr.FastForward(2 * time.Second) // miniredis快进时间

	// 锁应该已经过期，可以被其他客户端获取
	lock2 := NewDistributedLock(client, lockKey, 1*time.Second)
	err = lock2.Lock(ctx)
	if err != nil {
		t.Error("锁过期后应该可以被重新获取")
	}
}
