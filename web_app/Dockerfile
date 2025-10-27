# 多阶段构建Dockerfile
# 阶段1: 编译Go应用
FROM golang:1.24-alpine AS builder

# 设置工作目录
WORKDIR /build

# 安装必要工具
RUN apk add --no-cache git

# 复制go mod文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 生成Swagger文档
RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    swag init

# 编译应用（CGO_ENABLED=0 生成静态链接的二进制文件）
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bullbell-server .

# 阶段2: 运行环境
FROM alpine:latest

# 安装必要的运行时依赖
RUN apk --no-cache add ca-certificates tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

# 创建非root用户
RUN addgroup -g 1000 app && \
    adduser -D -u 1000 -G app app

# 设置工作目录
WORKDIR /app

# 从builder阶段复制编译好的二进制文件
COPY --from=builder /build/bullbell-server .

# 复制配置文件（使用docker专用配置）
COPY config.docker.yaml config.yaml

# 创建日志目录
RUN mkdir -p /app/logs && \
    chown -R app:app /app

# 切换到非root用户
USER app

# 暴露端口
EXPOSE 8082

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:8082/ping || exit 1

# 启动应用
CMD ["./bullbell-server"]


