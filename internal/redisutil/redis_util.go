package redisutil

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(client *redis.Client) *RedisClient {
	return &RedisClient{client}
}

// Set 设置键值对，支持任意类型
func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {

	// 将值序列化为 JSON 字符串
	val, err := json.Marshal(value)
	if err != nil {
		log.Printf("Failed to marshal value for key %s: %v", key, err)
		return err
	}

	err = r.client.Set(ctx, key, string(val), expiration).Err()
	if err != nil {
		log.Printf("Failed to set key %s: %v", key, err)
	}
	return err
}

// Get 获取键值对，支持任意类型
func (r *RedisClient) Get(ctx context.Context, key string, result interface{}) error {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			log.Printf("Key %s does not exist", key)
		} else {
			log.Printf("Failed to get key %s: %v", key, err)
		}
		return err
	}

	// 反序列化 JSON 字符串到指定类型
	err = json.Unmarshal([]byte(val), result)
	if err != nil {
		log.Printf("Failed to unmarshal value for key %s: %v", key, err)
		return err
	}

	return nil
}

// Del 删除键值对
func (r *RedisClient) Del(ctx context.Context, key string) error {
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		log.Printf("Failed to delete key %s: %v", key, err)
	}
	return err
}
