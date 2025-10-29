# Bullbell - 高性能论坛平台

[![Go CI/CD](https://github.com/DancingCircles/Blossom/actions/workflows/go.yml/badge.svg)](https://github.com/DancingCircles/Blossom/actions/workflows/go.yml)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://golang.org/)
[![GitHub stars](https://img.shields.io/github/stars/DancingCircles/Blossom?style=social)](https://github.com/DancingCircles/Blossom/stargazers)

> 一个生产级的 Go 论坛平台，支持高并发访问和现代化架构设计

## 项目简介

Bullbell 是一个现代化、高性能的论坛平台，专为可扩展性和可靠性而设计。采用 Go 语言和现代 Web 技术构建，展示了专业级的软件工程实践，包括清晰的分层架构、全面的测试覆盖和生产就绪的部署策略。

## 核心亮点

| 维度 | 说明 |
|------|------|
| **架构设计** | 三层架构（Controller-Logic-DAO），职责分离清晰 |
| **性能表现** | 轻量级接口 9,000+ QPS，数据库查询接口 3,800+ QPS |
| **技术栈** | Go 1.21+ / Gin / MySQL 8.0 / Redis 7 / Elasticsearch 8.11 |
| **安全机制** | JWT 认证、令牌桶限流、CORS 跨域保护 |
| **工程化** | Docker 容器化、CI/CD 流水线、自动化测试、API 文档 |
| **可扩展性** | Snowflake ID 生成、Redis 缓存、连接池、分布式友好设计 |

## 功能特性

### 用户管理
- 基于 JWT 的身份认证和授权
- bcrypt 密码加密
- 用户注册和登录，带完整验证

### 内容管理
- 话题创建、编辑和管理
- 层级评论系统
- 投票和评分机制
- 基于分类的内容组织

### 搜索与发现
- 基于 Elasticsearch 的全文搜索
- 搜索建议和自动补全
- 热门话题排行
- 分类统计信息

### 性能与可靠性
- Redis 缓存层，优化热点数据访问
- 令牌桶限流（每 IP 每秒 100 请求，突发容量 200）
- 结构化 JSON 日志
- Panic 恢复中间件
- 数据库连接池优化

## 界面展示

### 首页
![首页](web_app/pic/index.png)

### 话题列表
![话题列表](web_app/pic/show.png)

### 话题详情
![话题详情](web_app/pic/forum.png)

### 发布话题
![发布话题](web_app/pic/post.png)

## 技术架构

### 后端技术栈
- **开发语言**: Go 1.21+
- **Web 框架**: Gin（高性能 HTTP 框架）
- **数据库**: MySQL 8.0，带连接池优化
- **缓存**: Redis 7，用于会话和数据缓存
- **搜索引擎**: Elasticsearch 8.11，实现全文检索
- **API 文档**: Swagger/OpenAPI
- **ID 生成**: Snowflake 算法，支持分布式部署

### 前端技术栈
- **核心技术**: HTML5、CSS3、原生 JavaScript
- **设计风格**: 现代简约风 + 毛玻璃效果（Glassmorphism）
- **UI 组件**: 自定义响应式组件
- **构建工具**: 原生 ES 模块，无需打包工具

### 基础设施与 DevOps
- **容器化**: Docker 和 Docker Compose
- **CI/CD**: GitHub Actions 工作流
- **代码质量**: golangci-lint、gosec 安全扫描
- **测试**: 单元测试 + 竞态检测器
- **日志**: 结构化 JSON 日志，支持轮转

## 性能测试报告

使用 go-wrk 在本地开发环境测试结果：

| 端点类型 | QPS（请求/秒） | 平均响应时间 | P99 响应时间 | 错误率 |
|---------|---------------|-------------|-------------|--------|
| 健康检查 | 9,442 | 1.05ms | <1ms | 0% |
| 话题列表（数据库+缓存） | 3,853 | 5.24ms | ~1.64ms | 0% |
| 搜索（Elasticsearch） | 1,350 | 14.76ms | ~5.91ms | 0% |

**测试配置**: 100 并发连接，持续 10 秒，无限流

详细的压力测试报告请查看：[压力测试报告.md](压力测试报告.md)

## 快速开始

### 环境要求

- Go 1.21+
- MySQL 8.0
- Redis 7.0+
- Elasticsearch 8.11
- Docker & Docker Compose（推荐）

### 使用 Docker Compose 部署（推荐）

```bash
# 克隆项目
git clone https://github.com/DancingCircles/Blossom.git
cd Blossom

# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 访问应用
# 前端: http://localhost:80
# 后端: http://localhost:8082
# Swagger API 文档: http://localhost:8082/swagger/index.html
```

### 本地开发部署

#### 1. 启动依赖服务

```bash
# 启动 MySQL
docker run -d --name mysql \
  -e MYSQL_ROOT_PASSWORD=123456 \
  -e MYSQL_DATABASE=web_app \
  -p 13306:3306 \
  mysql:8.0

# 启动 Redis
docker run -d --name redis \
  -p 16379:6379 \
  redis:7-alpine

# 启动 Elasticsearch
docker run -d --name elasticsearch \
  -e "discovery.type=single-node" \
  -e "xpack.security.enabled=false" \
  -p 9200:9200 \
  elasticsearch:8.11.0
```

#### 2. 初始化数据库

```bash
cd web_app
mysql -h 127.0.0.1 -P 13306 -u root -p123456 web_app < sql/schema.sql
```

#### 3. 配置环境

编辑 `web_app/config.yaml` 文件，根据实际环境修改配置：

```yaml
app:
  name: "web_app"
  mode: "dev"
  port: 8082

mysql:
  host: "127.0.0.1"
  port: 13306
  user: "root"
  password: "123456"
  database: "web_app"

redis:
  host: "127.0.0.1"
  port: 16379

elasticsearch:
  url: "http://127.0.0.1:9200"
  index: "bullbell_topics"
```

#### 4. 启动后端服务

```bash
cd web_app

# 安装依赖
go mod download

# 生成 Swagger 文档
swag init

# 运行服务
go run main.go
```

#### 5. 启动前端服务

```bash
# 使用任意 HTTP 服务器
cd frontend

# 方式 1: Python
python -m http.server 3000

# 方式 2: Node.js
npx serve -p 3000
```

访问 http://localhost:3000 即可使用论坛。

## 项目结构

```
Bullbell/
├── .github/
│   └── workflows/
│       └── go.yml              # CI/CD 配置
├── frontend/                   # 前端代码
│   ├── css/                    # 样式文件
│   │   ├── style.css          # 主样式（含毛玻璃导航栏）
│   │   ├── auth.css           # 认证页样式
│   │   ├── post.css           # 发帖页样式
│   │   └── topic-search.css   # 搜索页样式
│   ├── js/                     # JavaScript 文件
│   │   ├── main.js            # 主逻辑（含导航栏自动隐藏）
│   │   ├── api.js             # API 封装
│   │   ├── auth.js            # 认证逻辑
│   │   ├── detail.js          # 详情页
│   │   └── post.js            # 发帖逻辑
│   ├── index.html              # 首页
│   ├── login.html              # 登录页
│   ├── post.html               # 发帖页
│   ├── detail.html             # 详情页
│   └── nginx.conf              # Nginx 配置
├── web_app/                    # 后端代码
│   ├── controllers/            # 控制器层
│   │   ├── admin.go           # 管理控制器
│   │   ├── comment.go         # 评论控制器
│   │   ├── search.go          # 搜索控制器
│   │   ├── topic.go           # 话题控制器
│   │   └── user.go            # 用户控制器
│   ├── dao/                    # 数据访问层
│   │   ├── mysql/             # MySQL 操作
│   │   ├── redis/             # Redis 操作
│   │   └── elasticsearch/     # Elasticsearch 操作
│   ├── logic/                  # 业务逻辑层
│   │   ├── comment.go         # 评论逻辑
│   │   ├── search.go          # 搜索逻辑
│   │   ├── topic.go           # 话题逻辑
│   │   └── user.go            # 用户逻辑
│   ├── models/                 # 数据模型
│   │   ├── comment.go         # 评论模型
│   │   ├── response.go        # 响应模型
│   │   ├── topic.go           # 话题模型
│   │   ├── user.go            # 用户模型
│   │   └── vote.go            # 投票模型
│   ├── middleware/             # 中间件
│   │   ├── cors.go            # CORS 跨域
│   │   ├── jwt.go             # JWT 认证
│   │   └── rate_limit.go      # 令牌桶限流
│   ├── routes/                 # 路由配置
│   │   └── routes.go          # 路由定义
│   ├── utils/                  # 工具函数
│   │   ├── jwt.go             # JWT 工具
│   │   ├── password.go        # 密码工具
│   │   └── snowflake.go       # ID 生成
│   ├── logger/                 # 日志系统
│   ├── settings/               # 配置管理
│   ├── docs/                   # Swagger 文档
│   ├── sql/                    # SQL 脚本
│   │   ├── schema.sql         # 数据库结构
│   │   └── migrate_to_snowflake.sql
│   ├── pic/                    # 项目截图
│   ├── config.yaml             # 配置文件
│   └── main.go                 # 入口文件
├── docs/                       # 文档目录
│   ├── CONTRIBUTING.md         # 贡献指南
│   ├── DOCKER_README.md        # Docker 说明
│   ├── IMPLEMENTATION_SUMMARY.md # 实现总结
│   └── FRONTEND_README.md      # 前端说明
├── docker-compose.yml          # Docker Compose 配置
├── Makefile                    # Make 命令
├── 压力测试报告.md              # 性能测试报告
└── README.md                   # 本文档
```

## 开发指南

### 代码规范

```bash
# 格式化代码
gofmt -w .

# 代码检查
go vet ./...

# 使用 golangci-lint
golangci-lint run
```

### 运行测试

```bash
cd web_app

# 运行所有测试
go test -v ./...

# 运行测试并生成覆盖率报告
go test -v -race -coverprofile=coverage.out ./...

# 查看覆盖率
go tool cover -html=coverage.out
```

### 压力测试

```bash
# 安装 go-wrk
go install github.com/tsliwowicz/go-wrk@latest

# 测试健康检查端点
go-wrk -c 100 -d 10 http://localhost:8082/ping

# 测试话题列表
go-wrk -c 50 -d 10 http://localhost:8082/api/v1/topics

# 测试搜索
go-wrk -c 20 -d 10 "http://localhost:8082/api/v1/search?q=test"
```

### 生成 Swagger 文档

```bash
cd web_app

# 安装 swag
go install github.com/swaggo/swag/cmd/swag@latest

# 生成文档
swag init

# 访问文档
# http://localhost:8082/swagger/index.html
```

## API 文档

启动服务后访问：
- **Swagger UI**: http://localhost:8082/swagger/index.html
- **API JSON**: http://localhost:8082/swagger/doc.json

主要 API 端点：

### 公开接口
- `POST /api/v1/register` - 用户注册
- `POST /api/v1/login` - 用户登录
- `GET /api/v1/topics` - 获取话题列表
- `GET /api/v1/topics/:id` - 获取话题详情
- `GET /api/v1/search` - 搜索话题
- `GET /api/v1/search/hot` - 热门话题

### 需要认证的接口
- `GET /api/v1/user/info` - 获取用户信息
- `POST /api/v1/topics` - 创建话题
- `POST /api/v1/topics/:id/vote` - 话题投票
- `POST /api/v1/topics/:id/comments` - 发表评论
- `DELETE /api/v1/comments/:id` - 删除评论

## 核心功能实现

### 令牌桶限流算法

使用 `golang.org/x/time/rate` 实现基于 IP 的令牌桶限流：

- 每个 IP 独立限流器
- 每秒生成 100 个令牌
- 桶容量 200（支持突发流量）
- 超限返回 429 状态码

详见：`web_app/middleware/rate_limit.go`

### Snowflake ID 生成

采用 Twitter Snowflake 算法生成全局唯一 ID：

- 64 位整数
- 支持分布式部署
- 时间有序
- 高性能（百万级 QPS）

详见：`web_app/utils/snowflake.go`

### Redis 缓存策略

- 话题列表缓存（TTL: 5 分钟）
- 热门话题缓存（TTL: 10 分钟）
- 用户会话缓存
- 缓存失效自动重建

详见：`web_app/dao/redis/`

### Elasticsearch 全文搜索

- 实时索引更新
- 分词搜索支持
- 高亮显示
- 搜索建议

详见：`web_app/dao/elasticsearch/`

## 前端特色

### 毛玻璃导航栏
- 半透明背景（70% 不透明度）
- backdrop-filter 毛玻璃效果
- 向下滚动自动隐藏
- 向上滚动自动显示
- 滚动到顶部完全透明

### 响应式设计
- 移动端适配
- 流畅的动画过渡
- 懒加载优化
- 骨架屏加载

## 部署建议

### 生产环境配置

1. **修改 CORS 配置**
```go
// web_app/middleware/cors.go
// 将 * 改为具体域名
c.Writer.Header().Set("Access-Control-Allow-Origin", "https://yourdomain.com")
```

2. **调整限流参数**
```go
// web_app/routes/routes.go
// 根据实际需求调整
middleware.IPRateLimit(500, 1000)
```

3. **启用 HTTPS**
- 配置 SSL 证书
- 强制 HTTPS 重定向

4. **数据库优化**
- 配置主从复制
- 启用慢查询日志
- 定期备份

5. **监控告警**
- 接入 Prometheus
- 配置 Grafana 面板
- 设置告警规则

## 贡献指南

欢迎贡献代码！请遵循以下流程：

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

详细的贡献指南请查看：[CONTRIBUTING.md](docs/CONTRIBUTING.md)

## 相关文档

- [实现总结](docs/IMPLEMENTATION_SUMMARY.md) - 详细的技术实现说明
- [Docker 部署](docs/DOCKER_README.md) - Docker 容器化部署指南
- [前端说明](docs/FRONTEND_README.md) - 前端技术详解
- [压力测试报告](压力测试报告.md) - 完整的性能测试数据

## 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

## 作者

项目维护者: [@DancingCircles](https://github.com/DancingCircles)

## 致谢

感谢所有为这个项目做出贡献的开发者！

---

如果这个项目对你有帮助，欢迎 Star 支持！
