.PHONY: help dev build clean test docker-up docker-down

# 默认目标
help:
	@echo "Blossom 项目管理命令"
	@echo ""
	@echo "使用方法: make [目标]"
	@echo ""
	@echo "可用目标:"
	@echo "  dev          - 启动开发环境"
	@echo "  build        - 构建后端应用"
	@echo "  clean        - 清理临时文件和构建产物"
	@echo "  test         - 运行测试"
	@echo "  test-cover   - 运行测试并生成覆盖率报告"
	@echo "  lint         - 运行代码检查"
	@echo "  fmt          - 格式化代码"
	@echo "  swagger      - 生成 Swagger 文档"
	@echo "  docker-up    - 启动所有 Docker 服务"
	@echo "  docker-down  - 停止所有 Docker 服务"
	@echo "  docker-logs  - 查看 Docker 日志"
	@echo "  install      - 安装开发依赖"

# 启动开发环境
dev:
	@echo "启动后端服务..."
	cd web_app && go run main.go

# 构建后端
build:
	@echo "构建后端应用..."
	cd web_app && go build -o ../bin/bullbell-server .
	@echo "构建完成: bin/bullbell-server"

# 清理临时文件
clean:
	@echo "清理临时文件..."
	rm -rf web_app/tmp/
	rm -rf web_app/logs/*.log
	rm -f web_app/web_app.exe
	rm -f web_app/coverage
	rm -f web_app/coverage.out
	rm -rf bin/
	@echo "清理完成！"

# 运行测试
test:
	@echo "运行测试..."
	cd web_app && go test -v ./...

# 测试覆盖率
test-cover:
	@echo "运行测试并生成覆盖率报告..."
	cd web_app && go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	cd web_app && go tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告: web_app/coverage.html"

# 代码检查
lint:
	@echo "运行代码检查..."
	cd web_app && go vet ./...
	cd web_app && gofmt -l .

# 格式化代码
fmt:
	@echo "格式化代码..."
	cd web_app && gofmt -w .
	@echo "代码格式化完成！"

# 生成 Swagger 文档
swagger:
	@echo "生成 Swagger 文档..."
	cd web_app && swag init
	@echo "Swagger 文档生成完成！"

# Docker 相关
docker-up:
	@echo "启动所有服务..."
	docker-compose up -d
	@echo "服务启动完成！"
	@echo "前端: http://localhost:8080"
	@echo "后端: http://localhost:8082"
	@echo "Swagger: http://localhost:8082/swagger/index.html"

docker-down:
	@echo "停止所有服务..."
	docker-compose down
	@echo "服务已停止！"

docker-logs:
	docker-compose logs -f

docker-restart:
	@echo "重启所有服务..."
	docker-compose restart
	@echo "服务已重启！"

# 安装开发依赖
install:
	@echo "安装开发依赖..."
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	cd web_app && go mod download
	@echo "依赖安装完成！"

# 初始化数据库
db-init:
	@echo "初始化数据库..."
	mysql -h 127.0.0.1 -P 13306 -u root -p123456 web_app < web_app/sql/schema.sql
	@echo "数据库初始化完成！"

# 同步数据到 Elasticsearch
sync-es:
	@echo "同步数据到 Elasticsearch..."
	curl -X POST http://localhost:8082/api/v1/admin/sync-es
	@echo "数据同步完成！"












