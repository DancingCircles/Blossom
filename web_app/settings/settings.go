// Package settings 负责应用程序配置管理
package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Conf 全局配置实例
var Conf = new(Config)

// Config 应用程序完整配置结构
type Config struct {
	App   *AppConfig   `mapstructure:"app"`
	MySQL *MysqlConfig `mapstructure:"mysql"`
	Redis *RedisConfig `mapstructure:"redis"`
	Log   *LogConfig   `mapstructure:"log"`
	Kafka *KafkaConfig `mapstructure:"kafka"`
}

// AppConfig 应用配置
type AppConfig struct {
	Name string `mapstructure:"name"`
	Mode string `mapstructure:"mode"`
	Port int    `mapstructure:"port"`
}

// MysqlConfig MySQL配置
type MysqlConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	Database     string `mapstructure:"database"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Database int    `mapstructure:"database"`
	PoolSize int    `mapstructure:"pool_size"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

// KafkaConfig Kafka配置
type KafkaConfig struct {
	Brokers []string `mapstructure:"brokers"`
	Topic   string   `mapstructure:"topic"`
	GroupID string   `mapstructure:"group_id"`
}

// Init 初始化配置系统
func Init() (err error) {
	// 设置配置文件
	viper.SetConfigFile("config.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// 读取配置文件
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("读取配置文件失败: %s \n", err)
		return err
	}

	// 启用配置热重载
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("检测到配置文件变化: %s\n", e.Name)
	})

	// 映射配置到结构体
	err = viper.Unmarshal(Conf)
	if err != nil {
		fmt.Printf("配置文件映射到结构体失败: %s\n", err)
		return err
	}

	// 打印配置信息
	fmt.Printf("配置文件路径: %s\n", viper.ConfigFileUsed())
	if Conf.MySQL != nil {
		fmt.Printf("MySQL配置 - 主机: %s, 端口: %d, 用户: %s, 数据库: %s\n",
			Conf.MySQL.Host, Conf.MySQL.Port, Conf.MySQL.User, Conf.MySQL.Database)
	}
	if Conf.App != nil {
		fmt.Printf("应用配置 - 名称: %s, 模式: %s, 端口: %d\n",
			Conf.App.Name, Conf.App.Mode, Conf.App.Port)
	}

	fmt.Println("配置系统初始化成功")
	return nil
}
