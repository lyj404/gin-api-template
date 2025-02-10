package bootstrap

import (
	"gin-api-template/config"

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
	app.Db = NewMySQLDataBase()

	// 按需初始化Redis
	if config.REnabled {
		app.Cache = InitRedis()
	}

	return *app
}

func (app *Application) CloseConnection() {
	CloseMySQLConnection(app.Db)
	CloseRedisConnection(app.Cache)
}
