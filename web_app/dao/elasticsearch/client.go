// Package elasticsearch 提供Elasticsearch连接和操作
package elasticsearch

import (
	"context"
	"fmt"

	"github.com/olivere/elastic/v7"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	client *elastic.Client
	index  string
)

// Init 初始化Elasticsearch客户端
func Init() error {
	var err error

	// 从配置读取ES地址
	esURL := viper.GetString("elasticsearch.url")
	index = viper.GetString("elasticsearch.index")
	sniff := viper.GetBool("elasticsearch.sniff")

	zap.L().Info("正在连接Elasticsearch", zap.String("url", esURL), zap.String("index", index))

	// 创建ES客户端
	client, err = elastic.NewClient(
		elastic.SetURL(esURL),
		elastic.SetSniff(sniff), // 单节点时设置为false
		elastic.SetHealthcheck(true),
		elastic.SetHealthcheckInterval(10),
	)
	if err != nil {
		return fmt.Errorf("连接Elasticsearch失败: %w", err)
	}

	// 检查连接
	ctx := context.Background()
	info, code, err := client.Ping(esURL).Do(ctx)
	if err != nil {
		return fmt.Errorf("Ping Elasticsearch失败: %w", err)
	}
	zap.L().Info("Elasticsearch连接成功",
		zap.String("cluster_name", info.ClusterName),
		zap.String("version", info.Version.Number),
		zap.Int("code", code))

	// 创建索引（如果不存在）
	err = createIndex(ctx)
	if err != nil {
		return fmt.Errorf("创建索引失败: %w", err)
	}

	return nil
}

// createIndex 创建话题索引
func createIndex(ctx context.Context) error {
	// 检查索引是否存在
	exists, err := client.IndexExists(index).Do(ctx)
	if err != nil {
		return err
	}

	if exists {
		zap.L().Info("索引已存在", zap.String("index", index))
		return nil
	}

	// 定义索引mapping
	mapping := `{
		"settings": {
			"number_of_shards": 1,
			"number_of_replicas": 0,
			"analysis": {
				"analyzer": {
					"ik_smart": {
						"type": "custom",
						"tokenizer": "standard"
					},
					"ik_max_word": {
						"type": "custom",
						"tokenizer": "standard"
					}
				}
			}
		},
		"mappings": {
			"properties": {
				"topic_id": {
					"type": "keyword"
				},
				"user_id": {
					"type": "keyword"
				},
				"title": {
					"type": "text",
					"analyzer": "standard",
					"search_analyzer": "standard",
					"fields": {
						"keyword": {
							"type": "keyword",
							"ignore_above": 256
						}
					}
				},
				"content": {
					"type": "text",
					"analyzer": "standard",
					"search_analyzer": "standard"
				},
				"category": {
					"type": "keyword"
				},
				"created_at": {
					"type": "date",
					"format": "yyyy-MM-dd HH:mm:ss||yyyy-MM-dd||epoch_millis"
				},
				"updated_at": {
					"type": "date",
					"format": "yyyy-MM-dd HH:mm:ss||yyyy-MM-dd||epoch_millis"
				},
				"view_count": {
					"type": "integer"
				},
				"comment_count": {
					"type": "integer"
				}
			}
		}
	}`

	// 创建索引
	createIndex, err := client.CreateIndex(index).BodyString(mapping).Do(ctx)
	if err != nil {
		return err
	}

	if !createIndex.Acknowledged {
		return fmt.Errorf("索引创建未被确认")
	}

	zap.L().Info("索引创建成功", zap.String("index", index))
	return nil
}

// GetClient 获取ES客户端
func GetClient() *elastic.Client {
	return client
}

// GetIndex 获取索引名称
func GetIndex() string {
	return index
}

// Close 关闭ES客户端
func Close() error {
	if client != nil {
		client.Stop()
		zap.L().Info("Elasticsearch连接已关闭")
	}
	return nil
}
