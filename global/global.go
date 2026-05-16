package global

import (
	"time"

	"github.com/lyj404/gin-api-template/internal/idgen"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	G_DB    *gorm.DB
	G_REDIS *redis.Client
)

type G_MODEL struct {
	ID        uint64         `gorm:"primarykey" json:"id,string"`   // 主键ID
	CreatedAt time.Time      `json:"created_at"`                    // 创建时间
	UpdatedAt time.Time      `json:"updated_at"`                    // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`                // 删除时间
}

func (m *G_MODEL) BeforeCreate(tx *gorm.DB) error {
	if m.ID == 0 {
		m.ID = idgen.NextID()
	}
	return nil
}
