// Package elasticsearch 提供搜索功能
package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
)

// SearchRequest 搜索请求参数
type SearchRequest struct {
	Keyword  string // 搜索关键词
	Category string // 分类筛选
	Page     int    // 页码
	PageSize int    // 每页数量
	SortBy   string // 排序字段: created_at, view_count, comment_count
}

// SearchResponse 搜索响应
type SearchResponse struct {
	Total    int64            `json:"total"`     // 总数
	Topics   []*TopicDocument `json:"topics"`    // 话题列表
	Page     int              `json:"page"`      // 当前页
	PageSize int              `json:"page_size"` // 每页数量
	Took     int64            `json:"took"`      // 耗时(毫秒)
}

// SearchTopics 搜索话题
func SearchTopics(req *SearchRequest) (*SearchResponse, error) {
	ctx := context.Background()

	// 构建查询
	boolQuery := elastic.NewBoolQuery()

	// 关键词搜索（标题和内容）
	if req.Keyword != "" {
		// 使用multi_match查询，支持中英文搜索
		multiMatch := elastic.NewMultiMatchQuery(req.Keyword, "title", "content").
			Type("best_fields").      // 最佳字段匹配
			TieBreaker(0.3).          // 多字段匹配时的权重
			Operator("OR").           // OR 操作符：任意一个词匹配即可
			MinimumShouldMatch("30%") // 至少匹配30%的词（对中文更友好）

		boolQuery = boolQuery.Must(multiMatch)
	}

	// 分类筛选
	if req.Category != "" {
		termQuery := elastic.NewTermQuery("category", req.Category)
		boolQuery = boolQuery.Filter(termQuery)
	}

	// 计算分页
	from := (req.Page - 1) * req.PageSize

	// 构建排序
	var sortBy string
	var sortOrder bool // true = 降序
	switch req.SortBy {
	case "view_count":
		sortBy = "view_count"
		sortOrder = false // 降序
	case "comment_count":
		sortBy = "comment_count"
		sortOrder = false // 降序
	case "created_at":
		fallthrough
	default:
		sortBy = "created_at"
		sortOrder = false // 最新的在前
	}

	// 执行搜索
	searchResult, err := client.Search().
		Index(index).
		Query(boolQuery).
		From(from).
		Size(req.PageSize).
		Sort(sortBy, sortOrder). // 排序
		Pretty(true).
		Do(ctx)

	if err != nil {
		zap.L().Error("搜索失败", zap.Error(err), zap.String("keyword", req.Keyword))
		return nil, err
	}

	// 解析结果
	topics := make([]*TopicDocument, 0)
	if searchResult.Hits != nil && searchResult.Hits.TotalHits.Value > 0 {
		for _, hit := range searchResult.Hits.Hits {
			var doc TopicDocument
			err := json.Unmarshal(hit.Source, &doc)
			if err != nil {
				zap.L().Error("解析搜索结果失败", zap.Error(err))
				continue
			}

			topics = append(topics, &doc)
		}
	}

	response := &SearchResponse{
		Total:    searchResult.Hits.TotalHits.Value,
		Topics:   topics,
		Page:     req.Page,
		PageSize: req.PageSize,
		Took:     searchResult.TookInMillis,
	}

	zap.L().Info("搜索完成",
		zap.String("keyword", req.Keyword),
		zap.Int64("total", response.Total),
		zap.Int("results", len(topics)),
		zap.Int64("took_ms", response.Took))

	return response, nil
}

// SuggestTopics 搜索建议（自动补全）
func SuggestTopics(prefix string, size int) ([]string, error) {
	ctx := context.Background()

	if prefix == "" {
		return []string{}, nil
	}

	// 使用前缀查询
	prefixQuery := elastic.NewPrefixQuery("title.keyword", prefix)

	searchResult, err := client.Search().
		Index(index).
		Query(prefixQuery).
		Size(size).
		FetchSourceContext(elastic.NewFetchSourceContext(true).Include("title")).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	suggestions := make([]string, 0)
	if searchResult.Hits != nil {
		for _, hit := range searchResult.Hits.Hits {
			var doc struct {
				Title string `json:"title"`
			}
			err := json.Unmarshal(hit.Source, &doc)
			if err == nil && doc.Title != "" {
				suggestions = append(suggestions, doc.Title)
			}
		}
	}

	return suggestions, nil
}

// GetTopicsByCategory 按分类获取热门话题
func GetTopicsByCategory(category string, size int) ([]*TopicDocument, error) {
	ctx := context.Background()

	termQuery := elastic.NewTermQuery("category", category)

	searchResult, err := client.Search().
		Index(index).
		Query(termQuery).
		Sort("view_count", false). // 按浏览量降序
		Size(size).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	topics := make([]*TopicDocument, 0)
	if searchResult.Hits != nil {
		for _, hit := range searchResult.Hits.Hits {
			var doc TopicDocument
			err := json.Unmarshal(hit.Source, &doc)
			if err == nil {
				topics = append(topics, &doc)
			}
		}
	}

	return topics, nil
}

// CountByCategory 统计各分类的话题数量
func CountByCategory() (map[string]int64, error) {
	ctx := context.Background()

	// 聚合查询
	agg := elastic.NewTermsAggregation().Field("category")

	searchResult, err := client.Search().
		Index(index).
		Size(0). // 不返回文档，只返回聚合结果
		Aggregation("categories", agg).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	result := make(map[string]int64)

	if agg, found := searchResult.Aggregations.Terms("categories"); found {
		for _, bucket := range agg.Buckets {
			if category, ok := bucket.Key.(string); ok {
				result[category] = bucket.DocCount
			}
		}
	}

	return result, nil
}

// DeleteIndex 删除索引（谨慎使用）
func DeleteIndex() error {
	ctx := context.Background()

	// 检查索引是否存在
	exists, err := client.IndexExists(index).Do(ctx)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("索引不存在: %s", index)
	}

	// 删除索引
	_, err = client.DeleteIndex(index).Do(ctx)
	if err != nil {
		return err
	}

	zap.L().Warn("索引已删除", zap.String("index", index))
	return nil
}
