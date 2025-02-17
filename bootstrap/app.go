package bootstrap

import (
	"gin-api-template/config"
	"gin-api-template/pkg/lib"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Application struct {
	Db    *gorm.DB
	Cache *redis.Client
}

func App() Application {
	app := &Application{}

	// 初始化数据库
	app.Db = lib.NewDataBase()

	// 按需初始化Redis
	if config.CfgRedis.Enabled {
		app.Cache = lib.InitRedis()
	}

	return *app
}

func (app *Application) CloseConnection() {
	lib.CloseMySQLConnection(app.Db)
	lib.CloseRedisConnection(app.Cache)
}
