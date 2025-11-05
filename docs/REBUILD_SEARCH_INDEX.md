# 重建搜索索引指南

## 问题说明
搜索功能已经升级，现在支持：
- ✅ 中文分词搜索（如：搜索"独立"可以找到"独立开发者"）
- ✅ 英文部分匹配（如：搜索"vvv"可以找到"VVVVVVVV"）
- ✅ 前缀搜索（支持输入前几个字符就能匹配）

## 重建索引步骤

### 1. 删除旧索引
```bash
curl -X DELETE "http://localhost:9200/topics"
```

### 2. 重启后端服务
后端服务会自动创建新的索引配置：
```bash
# Windows
cd web_app
go run main.go

# 或使用 Docker
docker-compose restart web_app
```

### 3. 同步数据到Elasticsearch
访问同步接口：
```bash
curl -X POST "http://localhost:8082/api/v1/admin/sync-es"
```

或者在浏览器中访问：
```
http://localhost:8082/api/v1/admin/sync-es
```

### 4. 验证搜索
访问前端页面，尝试搜索：
- 中文：`独立`、`开发`、`创业`
- 英文：`vvv`、`tech`、`design`

## 技术说明

### 新的分析器配置
- **索引时**：使用 `edge_ngram` 分词器（2-20字符），支持前缀匹配
- **搜索时**：使用 `standard` 分词器，标准分词
- **效果**：同时支持中文和英文的部分匹配搜索

### Edge N-gram 工作原理
例如文本 "VVVVVV" 会被分解为：
- "VV"
- "VVV" 
- "VVVV"
- "VVVVV"
- "VVVVVV"

这样搜索 "vvv" 就能匹配到 "VVVVVV"。

## 常见问题

### Q: 为什么要删除旧索引？
A: 因为索引的 mapping 和 analyzer 配置无法在线修改，必须重建索引。

### Q: 数据会丢失吗？
A: 不会。数据存储在 MySQL 中，Elasticsearch 只是搜索索引，重新同步即可。

### Q: 同步需要多久？
A: 取决于话题数量。通常100条话题不到1秒。

### Q: 如何验证索引是否创建成功？
A: 查看索引信息：
```bash
curl -X GET "http://localhost:9200/topics/_mapping?pretty"
```

应该能看到 `edge_ngram_tokenizer` 和 `index_analyzer` 配置。

