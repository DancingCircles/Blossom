# Blossom - 创意论坛社区

[![Go CI/CD](https://github.com/DancingCircles/Blossom/actions/workflows/go.yml/badge.svg)](https://github.com/DancingCircles/Blossom/actions/workflows/go.yml)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/DancingCircles/Blossom?style=social)](https://github.com/DancingCircles/Blossom/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/DancingCircles/Blossom?style=social)](https://github.com/DancingCircles/Blossom/network/members)

> 思想绽放的地方 | A place where ideas blossom

## 📖 项目简介

Blossom 是一个现代化的论坛社区平台，采用前后端分离架构，提供流畅的用户体验和强大的功能。

### ✨ 核心特性

- 🎨 **现代化UI设计** - 简约优雅的界面，出色的用户体验
- 🚀 **高性能架构** - Go后端 + Redis缓存 + Elasticsearch搜索
- 🔐 **安全可靠** - JWT认证、限流保护、数据验证
- 🔍 **智能搜索** - 基于Elasticsearch的全文搜索
- 💬 **实时交互** - 话题发布、评论、点赞等功能
- 📱 **响应式设计** - 完美适配各种设备

## 🏗️ 技术栈

### 后端
- **语言**: Go 1.21+
- **框架**: Gin Web Framework
- **数据库**: MySQL 8.0
- **缓存**: Redis 7
- **搜索**: Elasticsearch 8.11
- **文档**: Swagger/OpenAPI

### 前端
- **核心**: 原生 HTML5 + CSS3 + JavaScript
- **设计风格**: Neo-brutalism / Modern Minimalism
- **图标**: Emoji + SVG

### DevOps
- **容器化**: Docker + Docker Compose
- **CI/CD**: GitHub Actions
- **代码质量**: golangci-lint, gosec
- **测试**: Go testing + Race detector

## 🚀 快速开始

### 前置要求

- Go 1.21 或更高版本
- MySQL 8.0
- Redis 7
- Elasticsearch 8.11
- Docker (可选)

### 使用 Docker Compose (推荐)

```bash
# 克隆项目
git clone https://github.com/DancingCircles/Blossom.git
cd Blossom

# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 访问应用
# 前端: http://localhost:8080
# 后端: http://localhost:8082
# Swagger: http://localhost:8082/swagger/index.html
```

### 本地开发

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

```bash
# 复制配置文件
cp web_app/config.yaml web_app/config_local.yaml

# 编辑配置（如需要）
vim web_app/config_local.yaml
```

#### 4. 运行后端

```bash
cd web_app

# 安装依赖
go mod download

# 生成 Swagger 文档
swag init

# 运行服务
go run main.go
```

#### 5. 运行前端

```bash
# 使用任意 HTTP 服务器
cd frontend
python -m http.server 8080

# 或使用 Node.js
npx serve -p 8080
```

## 📁 项目结构

```
Bullbell/
├── .github/
│   └── workflows/
│       └── go.yml           # CI/CD 配置
├── frontend/                # 前端代码
│   ├── css/                 # 样式文件
│   ├── js/                  # JavaScript 文件
│   ├── index.html           # 主页
│   ├── login.html           # 登录页
│   ├── post.html            # 发帖页
│   └── detail.html          # 详情页
├── web_app/                 # 后端代码
│   ├── controllers/         # 控制器层
│   ├── dao/                 # 数据访问层
│   │   ├── mysql/          # MySQL
│   │   ├── redis/          # Redis
│   │   └── elasticsearch/  # Elasticsearch
│   ├── logic/              # 业务逻辑层
│   ├── models/             # 数据模型
│   ├── middleware/         # 中间件
│   ├── routes/             # 路由配置
│   ├── utils/              # 工具函数
│   ├── logger/             # 日志系统
│   ├── settings/           # 配置管理
│   ├── docs/               # Swagger 文档
│   └── main.go             # 入口文件
├── docker-compose.yml       # Docker Compose 配置
├── .gitignore              # Git 忽略规则
└── README.md               # 项目文档
```

## 🔧 开发指南

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

### 生成 Swagger 文档

```bash
cd web_app

# 安装 swag
go install github.com/swaggo/swag/cmd/swag@latest

# 生成文档
swag init

# 访问 http://localhost:8082/swagger/index.html
```

## 📊 API 文档

启动服务后访问：
- Swagger UI: http://localhost:8082/swagger/index.html
- API JSON: http://localhost:8082/swagger/doc.json

## 🤝 贡献指南

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 📝 开发日志

- 详细的实现总结请查看 [IMPLEMENTATION_SUMMARY.md](docs/IMPLEMENTATION_SUMMARY.md)
- Docker 部署说明请查看 [DOCKER_README.md](docs/DOCKER_README.md)
- 前端动态加载说明请查看 [DYNAMIC_LOADING_README.md](docs/DYNAMIC_LOADING_README.md)
- 贡献指南请查看 [CONTRIBUTING.md](docs/CONTRIBUTING.md)

## 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

## 👥 作者

- 项目维护者: [@DancingCircles](https://github.com/DancingCircles)

## 🙏 致谢

感谢所有为这个项目做出贡献的开发者！

---

⭐ 如果这个项目对你有帮助，请给我们一个星标！

