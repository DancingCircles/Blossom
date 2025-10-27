// Package utils_test 提供工具函数测试
package utils

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestInitSnowflake 测试雪花算法初始化
func TestInitSnowflake(t *testing.T) {
	err := InitSnowflake()
	assert.NoError(t, err, "雪花算法初始化应该成功")
}

// TestGenerateID 测试ID生成
func TestGenerateID(t *testing.T) {
	// 确保已初始化
	err := InitSnowflake()
	assert.NoError(t, err)

	// 生成ID
	id := GenerateID()
	assert.Greater(t, id, int64(0), "生成的ID应该大于0")
}

// TestGenerateUniqueIDs 测试ID唯一性
func TestGenerateUniqueIDs(t *testing.T) {
	// 确保已初始化
	err := InitSnowflake()
	assert.NoError(t, err)

	// 生成1000个ID
	ids := make(map[int64]bool)
	for i := 0; i < 1000; i++ {
		id := GenerateID()

		// 检查是否重复
		if ids[id] {
			t.Errorf("生成了重复的ID: %d", id)
		}
		ids[id] = true
	}

	// 验证生成了1000个唯一ID
	assert.Equal(t, 1000, len(ids), "应该生成1000个唯一ID")
}

// TestGenerateIDMonotonic 测试ID单调递增性
func TestGenerateIDMonotonic(t *testing.T) {
	// 确保已初始化
	err := InitSnowflake()
	assert.NoError(t, err)

	// 生成100个ID，验证递增
	prevID := GenerateID()
	for i := 0; i < 100; i++ {
		currentID := GenerateID()
		assert.Greater(t, currentID, prevID, "ID应该单调递增")
		prevID = currentID
	}
}

// TestGenerateIDConcurrent 测试并发生成ID
func TestGenerateIDConcurrent(t *testing.T) {
	// 确保已初始化
	err := InitSnowflake()
	assert.NoError(t, err)

	// 并发生成10000个ID
	const goroutines = 100
	const idsPerGoroutine = 100
	totalIDs := goroutines * idsPerGoroutine

	ids := make([]int64, totalIDs)
	var wg sync.WaitGroup

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()
			for j := 0; j < idsPerGoroutine; j++ {
				ids[start+j] = GenerateID()
			}
		}(i * idsPerGoroutine)
	}

	wg.Wait()

	// 验证所有ID唯一
	idMap := make(map[int64]bool)
	for _, id := range ids {
		if idMap[id] {
			t.Errorf("并发生成了重复的ID: %d", id)
		}
		idMap[id] = true
		assert.Greater(t, id, int64(0), "ID应该大于0")
	}

	assert.Equal(t, totalIDs, len(idMap), "所有ID应该唯一")
}

// TestGenerateIDFormat 测试ID格式（19位数字）
func TestGenerateIDFormat(t *testing.T) {
	// 确保已初始化
	err := InitSnowflake()
	assert.NoError(t, err)

	// 生成10个ID
	for i := 0; i < 10; i++ {
		id := GenerateID()

		// 雪花算法ID应该是19位数字（最大值：9223372036854775807）
		assert.Greater(t, id, int64(0), "ID应该大于0")
		assert.LessOrEqual(t, id, int64(9223372036854775807), "ID应该在int64范围内")
	}
}

// TestGenerateIDPerformance 测试ID生成性能
func TestGenerateIDPerformance(t *testing.T) {
	// 确保已初始化
	err := InitSnowflake()
	assert.NoError(t, err)

	// 生成100000个ID，验证性能
	ids := make([]int64, 100000)
	for i := 0; i < 100000; i++ {
		ids[i] = GenerateID()
	}

	// 验证所有ID唯一
	idMap := make(map[int64]bool)
	for _, id := range ids {
		assert.False(t, idMap[id], "不应该有重复ID")
		idMap[id] = true
	}

	assert.Equal(t, 100000, len(idMap), "应该生成100000个唯一ID")
}

// BenchmarkGenerateID 基准测试：ID生成性能
func BenchmarkGenerateID(b *testing.B) {
	// 确保已初始化
	if err := InitSnowflake(); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GenerateID()
	}
}

// BenchmarkGenerateIDParallel 基准测试：并发ID生成性能
func BenchmarkGenerateIDParallel(b *testing.B) {
	// 确保已初始化
	if err := InitSnowflake(); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			GenerateID()
		}
	})
}
