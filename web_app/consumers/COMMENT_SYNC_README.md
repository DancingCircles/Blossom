# 评论同步逻辑说明

## 功能概述

当评论表（`comments`）发生变更时（INSERT/UPDATE/DELETE），自动同步话题的评论数到 Elasticsearch。

## 实现原理

```
用户评论 → MySQL comments表 
    ↓
Canal 监听 binlog
    ↓
发送到 Kafka (canal-topic)
    ↓
ES Consumer 消费消息
    ↓
1. 提取 topic_id
2. 查询 MySQL 最新评论数
3. 更新 ES 中的话题文档
```

## 核心代码

### 1. ES 更新方法 (elasticsearch/topic.go)

```go
// UpdateTopicCommentCount 更新话题的评论数
func UpdateTopicCommentCount(topicID int64, commentCount int) error
```

### 2. 消费者处理方法 (consumers/es_consumer.go)

```go
// handleCommentChange 处理评论变更（更新话题的评论数）
func (c *ESConsumer) handleCommentChange(msg *CanalMessage) error
```

## 特性

✅ **去重处理**：一次批量变更中，相同话题只更新一次  
✅ **容错机制**：单个话题失败不影响其他话题的更新  
✅ **日志记录**：详细记录每次同步的结果  
✅ **幂等性**：多次更新结果一致  
✅ **文档不存在处理**：ES 中不存在的话题会记录警告但不报错  

## 数据流示例

### 场景1：用户发表评论

1. 用户通过 API 发表评论 → MySQL `comments` 表插入记录
2. Canal 监听到 INSERT 操作 → 发送到 Kafka
3. ES Consumer 接收消息：
   ```json
   {
     "type": "INSERT",
     "table": "comments",
     "data": [{"topic_id": 12345, ...}]
   }
   ```
4. 提取 `topic_id = 12345`
5. 查询 MySQL: `SELECT COUNT(*) FROM comments WHERE topic_id = 12345` → 结果: 5
6. 更新 ES: `PUT /bullbell_topics/_doc/12345 { "comment_count": 5 }`
7. 日志输出：
   ```
   INFO 话题评论数已同步到ES {"topic_id": 12345, "comment_count": 5, "operation": "INSERT"}
   ```

### 场景2：用户删除评论

1. 用户删除评论 → MySQL `comments` 表删除记录
2. Canal 监听到 DELETE 操作 → 发送到 Kafka
3. ES Consumer 处理流程同上
4. 更新后的评论数会自动减少

### 场景3：批量操作

如果一次有多条评论变更：
```json
{
  "data": [
    {"topic_id": 100, ...},
    {"topic_id": 100, ...},  // 重复的会被去重
    {"topic_id": 200, ...}
  ]
}
```

只会执行2次更新（topic 100 和 topic 200）

## 监控和调试

### 查看日志

```bash
# 查看同步日志
tail -f web_app/logs/web_app.log | grep "话题评论数"
```

### 验证同步结果

```bash
# 1. 查询 MySQL 评论数
mysql> SELECT topic_id, COUNT(*) as count FROM comments GROUP BY topic_id;

# 2. 查询 ES 评论数
curl -X GET "localhost:9200/bullbell_topics/_search" -H 'Content-Type: application/json' -d'
{
  "query": {"match_all": {}},
  "_source": ["topic_id", "title", "comment_count"]
}'
```

### 手动触发同步

如果发现数据不一致，可以：

1. **单个话题同步**：发表或删除一条评论
2. **全量同步**：调用管理接口
   ```bash
   curl -X POST "localhost:8082/api/v1/admin/sync-es"
   ```

## 错误处理

### 常见错误

1. **MySQL 查询失败**
   ```
   ERROR 查询话题评论数失败 {"error": "...", "topic_id": 12345}
   ```
   - 原因：MySQL 连接断开或话题不存在
   - 处理：记录错误，继续处理其他话题

2. **ES 更新失败**
   ```
   ERROR 更新ES话题评论数失败 {"error": "...", "topic_id": 12345}
   ```
   - 原因：ES 连接问题或索引不存在
   - 处理：记录错误，继续处理其他话题

3. **话题索引不存在**
   ```
   WARN 话题索引不存在，跳过更新评论数 {"topic_id": "12345"}
   ```
   - 原因：话题还未同步到 ES
   - 处理：记录警告，不报错（下次话题同步时会创建索引）

## 性能考虑

- **批量去重**：使用 map 去重，避免重复查询
- **继续执行**：单个失败不影响其他话题
- **异步处理**：通过 Kafka 异步消费，不阻塞主流程
- **索引刷新**：使用 `refresh=true` 确保数据立即可见

## 未来优化方向

- [ ] 批量更新 ES（如果同时有多个话题）
- [ ] 增加重试机制（失败时自动重试）
- [ ] 增加 Prometheus 指标监控
- [ ] 支持增量更新（+1/-1 而不是查询总数）

