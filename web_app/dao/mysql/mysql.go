// Package mysql 提供MySQL数据库连接管理
package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"web_app/models"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// db 全局数据库连接池实例
var db *sqlx.DB

// Init 初始化MySQL数据库连接
func Init() (err error) {
	// 读取数据库配置
	host := viper.GetString("mysql.host")
	port := viper.GetInt("mysql.port")
	user := viper.GetString("mysql.user")
	password := viper.GetString("mysql.password")
	database := viper.GetString("mysql.database")

	// 设置默认值
	if host == "" {
		host = "127.0.0.1"
	}
	if port == 0 {
		port = 13306
	}
	if user == "" {
		return fmt.Errorf("MySQL用户名不能为空")
	}
	if database == "" {
		return fmt.Errorf("MySQL数据库名不能为空")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, database)

	// 连接数据库
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Error("MySQL 数据库连接失败", zap.Error(err))
		return err
	}

	// 配置连接池
	db.SetMaxOpenConns(viper.GetInt("mysql.max_open_conns"))
	db.SetMaxIdleConns(viper.GetInt("mysql.max_idle_conns"))
	zap.L().Info("MySQL 数据库连接成功",
		zap.String("host", viper.GetString("mysql.host")),
		zap.Int("port", viper.GetInt("mysql.port")),
		zap.String("database", viper.GetString("mysql.database")),
		zap.Int("max_open_conns", viper.GetInt("mysql.max_open_conns")),
		zap.Int("max_idle_conns", viper.GetInt("mysql.max_idle_conns")),
	)

	return nil
}

// Close 关闭数据库连接
func Close() {
	if db != nil {
		if err := db.Close(); err != nil {
			zap.L().Error("关闭 MySQL 连接失败", zap.Error(err))
		} else {
			zap.L().Info("MySQL 连接已关闭")
		}
	}
}

// GetDB 获取数据库连接实例
func GetDB() *sqlx.DB {
	return db
}

func CheckUserExists(username string) (bool, error) {
	var user models.User
	err := db.Get(&user, "SELECT * FROM users WHERE username = ?", username)
	if err != nil {
		// 如果是"无结果"错误，说明用户不存在，返回false
		if err == sql.ErrNoRows {
			return false, nil
		}
		// 其他错误直接返回
		return false, err
	}
	// 查询成功且有结果，说明用户存在
	return true, nil
}

func InsertUser(userID int64, username, email, hashedPassword string) error {
	stmt := "INSERT INTO users (id, username, email, password) VALUES (?, ?, ?, ?)"
	if _, err := db.Exec(stmt, userID, username, email, hashedPassword); err != nil {
		return errors.New("插入用户失败: " + err.Error())
	}
	return nil
}

// GetUserByUsername 根据用户名获取用户信息
func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := db.Get(&user, "SELECT id, username, email, password, created_at, updated_at FROM users WHERE username = ?", username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByID 根据用户ID获取用户信息
func GetUserByID(userID int64) (*models.User, error) {
	var user models.User
	err := db.Get(&user, "SELECT id, username, email, password, created_at, updated_at FROM users WHERE id = ?", userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return &user, nil
}
