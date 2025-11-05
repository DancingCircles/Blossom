# Goweb-Frame

一个基于 Gin + MySQL + Redis 的 Go Web 脚手架

## 🚀 特性

- 🏗️ **清晰架构**: 分层设计，职责分明
- ⚡ **高性能**: Gin框架 + 连接池优化
- 📝 **结构化日志**: Zap日志系统，支持切割
- 🗄️ **双存储**: MySQL + Redis
- ⚙️ **配置管理**: Viper热重载
- 🔧 **生产就绪**: 优雅关机、错误恢复

## 📋 技术栈

- **Web框架**: Gin v1.10.1
- **数据库**: MySQL + SQLX
- **缓存**: Redis + go-redis
- **配置**: Viper
- **日志**: Zap + Lumberjack

## 🏗️ 项目结构

```
web_app/
├── main.go              # 程序入口
├── config.yaml          # 配置文件
├── settings/            # 配置管理
├── logger/              # 日志系统
├── dao/                 # 数据访问层
│   ├── mysql/          # MySQL操作
│   └── redis/          # Redis操作
├── routes/             # 路由配置
├── controllers/        # 控制器（待扩展）
├── logic/             # 业务逻辑（待扩展）
├── models/            # 数据模型（待扩展）
└── pkg/               # 工具包（待扩展）
```

## 🚀 快速开始

### 1. 环境要求
- Go 1.24.5+
- MySQL 5.7+
- Redis 6.0+

### 2. 启动项目
```bash
git clone https://github.com/DancingCircles/Goweb-Frame.git
cd Goweb-Frame
go mod tidy

# 修改 config.yaml 中的数据库配置
go run main.go
```

### 3. 测试
```bash
curl http://localhost:8081/ping
# 返回: {"message":"pong","status":"healthy","service":"Goweb-Frame"}
```

## ⚙️ 配置

编辑 `config.yaml`：
```yaml
app:
  port: 8081

mysql:
  host: "127.0.0.1"
  user: "root"
  password: "123456"
  database: "web_app"

redis:
  host: "127.0.0.1"
  password: "123456"
```

## 📝 添加接口

### 1. 添加路由
```go
// routes/routes.go
r.GET("/api/users", controllers.GetUsers)
```

### 2. 创建控制器
```go
// controllers/user.go
func GetUsers(c *gin.Context) {
    c.JSON(200, gin.H{"data": "用户列表"})
}
```

### 3. 数据库操作
```go
db := mysql.GetDB()
rows, err := db.Query("SELECT * FROM users")
```

### 4. Redis操作
```go
rdb := redis.GetRedis()
rdb.Set(context.Background(), "key", "value", 0)
```

## 🛠️ 部署

### 开发环境
```bash
go run main.go
```

### 生产环境
```bash
go build -o web_app main.go
./web_app
```

## 📄 许可证

MIT License

---

⭐ 如果这个项目对你有帮助，请给个 Star！