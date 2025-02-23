package global

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	G_DB    *gorm.DB
	G_REDIS *redis.Client
)
