# 贡献指南

感谢你对 Blossom 项目的关注！这份文档将帮助你快速上手项目开发。

## 🚀 开发环境设置

### 1. 环境要求

- Go 1.21+
- MySQL 8.0
- Redis 7
- Elasticsearch 8.11
- Git
- Make (可选，用于快捷命令)

### 2. 克隆项目

```bash
git clone https://github.com/YOUR_USERNAME/Bullbell.git
cd Bullbell
```

### 3. 安装依赖

```bash
# 使用 Makefile
make install

# 或手动安装
cd web_app
go mod download
go install github.com/swaggo/swag/cmd/swag@latest
```

### 4. 启动服务

```bash
# 使用 Docker Compose（推荐）
docker-compose up -d

# 或手动启动
make dev
```

## 📝 开发规范

### 代码风格

1. **Go 代码**
   - 遵循 [Effective Go](https://golang.org/doc/effective_go.html)
   - 使用 `gofmt` 格式化代码
   - 运行 `go vet` 进行静态检查

```bash
# 格式化代码
make fmt

# 代码检查
make lint
```

2. **前端代码**
   - 使用 2 空格缩进
   - 注释清晰，易于理解
   - 保持代码简洁

### 提交规范

使用语义化的提交信息：

```
<类型>(<范围>): <描述>

[可选的正文]

[可选的脚注]
```

**类型：**
- `feat`: 新功能
- `fix`: 修复 bug
- `docs`: 文档更新
- `style`: 代码格式调整（不影响功能）
- `refactor`: 代码重构
- `perf`: 性能优化
- `test`: 测试相关
- `chore`: 构建/工具链更新

**示例：**
```
feat(search): 添加 Elasticsearch 全文搜索功能

- 实现关键词搜索
- 支持模糊匹配
- 添加搜索结果高亮

Closes #123
```

### 分支策略

- `main`: 主分支，保持稳定
- `develop`: 开发分支
- `feature/*`: 功能分支
- `fix/*`: 修复分支
- `hotfix/*`: 紧急修复分支

## 🔄 开发流程

### 1. 创建功能分支

```bash
git checkout -b feature/your-feature-name
```

### 2. 开发与测试

```bash
# 运行测试
make test

# 查看覆盖率
make test-cover
```

### 3. 提交代码

```bash
git add .
git commit -m "feat: 你的功能描述"
```

### 4. 推送分支

```bash
git push origin feature/your-feature-name
```

### 5. 创建 Pull Request

在 GitHub 上创建 PR，并：
- 填写清晰的描述
- 关联相关 Issue
- 等待 CI 检查通过
- 等待代码审查

## ✅ 测试要求

### 单元测试

- 为新功能编写测试
- 保持测试覆盖率 > 70%
- 测试要快速、稳定、独立

```bash
# 运行测试
cd web_app
go test -v ./...

# 测试特定包
go test -v ./logic

# 带覆盖率
go test -v -coverprofile=coverage.out ./...
```

### 集成测试

在提交 PR 前：
- 手动测试核心功能
- 确保 API 正常工作
- 验证前端交互

## 🐛 Bug 报告

提交 Bug 时请包含：

1. **环境信息**
   - 操作系统
   - Go 版本
   - 浏览器（前端问题）

2. **重现步骤**
   - 详细的操作步骤
   - 期望结果
   - 实际结果

3. **相关信息**
   - 错误日志
   - 截图（如适用）
   - 相关代码片段

## 💡 功能建议

提交功能建议时请说明：

1. **使用场景** - 为什么需要这个功能
2. **期望行为** - 功能应该如何工作
3. **替代方案** - 是否考虑过其他方案
4. **影响范围** - 可能影响哪些模块

## 📚 学习资源

### Go 相关
- [Go 官方文档](https://golang.org/doc/)
- [Gin 框架文档](https://gin-gonic.com/docs/)
- [Go 设计模式](https://github.com/tmrts/go-patterns)

### 前端相关
- [MDN Web Docs](https://developer.mozilla.org/)
- [CSS Tricks](https://css-tricks.com/)

### 数据库相关
- [MySQL 文档](https://dev.mysql.com/doc/)
- [Redis 文档](https://redis.io/documentation)
- [Elasticsearch 指南](https://www.elastic.co/guide/)

## 🙋 获取帮助

如果遇到问题：

1. 查看[项目文档](README.md)
2. 搜索[已有 Issues](https://github.com/YOUR_USERNAME/Bullbell/issues)
3. 在 Discord/微信群提问
4. 创建新的 Issue

## 📄 许可证

通过贡献代码，你同意你的贡献将在 MIT 许可证下发布。

---

再次感谢你的贡献！🎉

