package global

import (
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	G_DB    *gorm.DB
	G_REDIS *redis.Client
)

type G_MODEL struct {
	ID        uint `gorm:"primarykey"` // 主键ID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 删除时间
}
