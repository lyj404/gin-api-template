package bootstrap

import (
	"context"
	"gin-api-template/config"
	"log"

	"github.com/redis/go-redis/v9"
)

func InitRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     config.RHost + ":" + config.RPort,
		Password: config.RPassWord,
		DB:       config.RDataBase,
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
