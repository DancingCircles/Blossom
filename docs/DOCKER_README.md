# Docker 部署指南

## 快速启动

### 1. 启动所有服务
```bash
docker-compose up -d
```

### 2. 查看服务状态
```bash
docker-compose ps
```

### 3. 查看日志
```bash
# 查看所有服务日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f backend
```

### 4. 停止服务
```bash
docker-compose stop
```

### 5. 重启服务
```bash
docker-compose restart
```

### 6. 完全删除（包括数据）
```bash
docker-compose down -v
```

## 访问地址

- **前端页面**: http://localhost
- **API接口**: http://localhost:8082
- **Swagger文档**: http://localhost:8082/swagger/index.html
- **MySQL**: localhost:3306
- **Redis**: localhost:6379

## 服务架构

```
┌─────────────┐
│  Frontend   │ :80
│   (Nginx)   │
└──────┬──────┘
       │
       ▼
┌─────────────┐     ┌─────────┐     ┌─────────┐
│   Backend   │────▶│  MySQL  │     │  Redis  │
│    (Go)     │ :8082│         │:3306│         │:6379
└─────────────┘     └─────────┘     └─────────┘
```

## 首次启动注意事项

1. 首次启动会自动初始化数据库（约30秒）
2. 等待所有服务健康检查通过
3. 访问 http://localhost/swagger/index.html 查看API文档

## 数据持久化

数据保存在Docker volume中：
- `mysql_data` - MySQL数据
- `redis_data` - Redis数据

删除volume会清空所有数据，请谨慎操作！


