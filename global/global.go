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
	ID        uint           `gorm:"primarykey" json:"id"`          // 主键ID
	CreatedAt time.Time      `json:"created_at"`                    // 创建时间
	UpdatedAt time.Time      `json:"updated_at"`                    // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`                // 删除时间
}
