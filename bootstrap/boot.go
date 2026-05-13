package bootstrap

import (
	"github.com/lyj404/gin-api-template/config"
	"github.com/lyj404/gin-api-template/global"
	"github.com/lyj404/gin-api-template/pkg/lib"
)

func Boot() {
	// 初始化数据库
	global.G_DB = lib.NewDataBase()

	// 按需初始化Redis
	if config.CfgRedis.Enabled {
		global.G_REDIS = lib.InitRedis()
	}

}

// BootDBOnly 仅初始化数据库，适用于 CLI 等不依赖 Redis 的场景
func BootDBOnly() {
	global.G_DB = lib.NewDataBase()
}

func CloseConnection() {
	lib.CloseDataBaseConnection(global.G_DB)
	lib.CloseRedisConnection(global.G_REDIS)
}
