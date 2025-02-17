package lib

import (
	"context"
	"gin-api-template/config"
	"log"

	"github.com/redis/go-redis/v9"
)

func InitRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     config.CfgRedis.Host + ":" + config.CfgRedis.Port,
		Password: config.CfgRedis.Password,
		DB:       config.CfgRedis.Database,
	})
	ping := client.Ping(context.Background())
	err := ping.Err()
	if err != nil {
		log.Fatal("failed to connect to Redis:", err)
	}
	return client
}

func CloseRedisConnection(client *redis.Client) {
	if client == nil {
		return
	}

	// 关闭Redis连接
	err := client.Close()
	if err != nil {
		log.Fatal("Redis连接关闭失败:", err)
	}
}
