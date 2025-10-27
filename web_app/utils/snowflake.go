// Package utils 提供通用工具函数
package utils

import (
	"fmt"
	"sync"
	"time"

	"github.com/spf13/viper"
)

const (
	// 机器ID占用位数
	machineIDBits = uint(10)
	// 序列号占用位数
	sequenceBits = uint(12)

	// 最大机器ID (1023)
	maxMachineID = int64(-1) ^ (int64(-1) << machineIDBits)
	// 最大序列号 (4095)
	maxSequence = int64(-1) ^ (int64(-1) << sequenceBits)

	// 时间戳左移位数
	timestampShift = machineIDBits + sequenceBits
	// 机器ID左移位数
	machineIDShift = sequenceBits

	// 自定义起始时间戳 (2024-01-01 00:00:00 UTC)
	epoch = int64(1704067200000) // 毫秒
)

// Snowflake 雪花算法ID生成器
type Snowflake struct {
	mu        sync.Mutex // 互斥锁
	timestamp int64      // 上次生成ID的时间戳
	machineID int64      // 机器ID
	sequence  int64      // 序列号
}

var (
	snowflakeInstance *Snowflake
	once              sync.Once
)

// InitSnowflake 初始化雪花算法生成器
func InitSnowflake() error {
	var initErr error
	once.Do(func() {
		// 从配置文件读取机器ID，默认为1
		machineID := int64(viper.GetInt("snowflake.machine_id"))
		if machineID == 0 {
			machineID = 1
		}

		// 验证机器ID范围
		if machineID < 0 || machineID > maxMachineID {
			initErr = fmt.Errorf("机器ID必须在 0 到 %d 之间", maxMachineID)
			return
		}

		snowflakeInstance = &Snowflake{
			timestamp: 0,
			machineID: machineID,
			sequence:  0,
		}
	})
	return initErr
}

// GetSnowflake 获取雪花算法生成器实例
func GetSnowflake() *Snowflake {
	if snowflakeInstance == nil {
		panic("雪花算法生成器未初始化，请先调用 InitSnowflake()")
	}
	return snowflakeInstance
}

// NextID 生成下一个唯一ID
func (s *Snowflake) NextID() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 获取当前时间戳（毫秒）
	now := time.Now().UnixMilli()

	// 如果当前时间小于上次生成ID的时间，说明系统时钟回拨，抛出错误
	if now < s.timestamp {
		panic(fmt.Sprintf("时钟回拨detected，拒绝生成ID。上次时间: %d, 当前时间: %d", s.timestamp, now))
	}

	// 如果是同一毫秒内生成ID
	if now == s.timestamp {
		// 序列号自增
		s.sequence = (s.sequence + 1) & maxSequence
		// 如果序列号溢出（达到4096），则等待下一毫秒
		if s.sequence == 0 {
			now = s.waitNextMillis(now)
		}
	} else {
		// 不同毫秒，序列号重置为0
		s.sequence = 0
	}

	// 更新时间戳
	s.timestamp = now

	// 生成ID：
	// 1. 时间戳部分：(now - epoch) << timestampShift
	// 2. 机器ID部分：machineID << machineIDShift
	// 3. 序列号部分：sequence
	id := ((now - epoch) << timestampShift) |
		(s.machineID << machineIDShift) |
		s.sequence

	return id
}

// waitNextMillis 等待下一毫秒
func (s *Snowflake) waitNextMillis(currentMillis int64) int64 {
	for currentMillis == s.timestamp {
		currentMillis = time.Now().UnixMilli()
	}
	return currentMillis
}

// GenerateID 全局函数：生成雪花算法ID
func GenerateID() int64 {
	return GetSnowflake().NextID()
}
