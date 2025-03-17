package bootstrap

import (
	"gin-api-template/config"
	"gin-api-template/global"
	"gin-api-template/pkg/lib"
)

func Boot() {
	// 初始化配置
	config.InitConfig()

	// 初始化数据库
	global.G_DB = lib.NewDataBase()

	// 按需初始化Redis
	if config.CfgRedis.Enabled {
		global.G_REDIS = lib.InitRedis()
	}

}

func CloseConnection() {
	lib.CloseMySQLConnection(global.G_DB)
	lib.CloseRedisConnection(global.G_REDIS)
}
