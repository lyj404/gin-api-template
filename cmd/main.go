package main

import (
	"github.com/lyj404/gin-api-template/bootstrap"
	"github.com/lyj404/gin-api-template/config"
)

func main() {
	// 初始化配置
	config.InitConfig()

	// 初始化数据库和缓存（用于设置全局变量）
	bootstrap.Boot()

	// 使用 Wire 初始化应用
	app, err := InitializeApp()
	if err != nil {
		panic(err)
	}

	// 程序结束时关闭资源
	defer func() {
		if app.Logger != nil {
			app.Logger.Sync()
		}
		bootstrap.CloseConnection()
	}()

	// 设置swagger文档
	config.SetUpSwag()

	// 注册所有路由
	app.RegisterRoutes()

	// 运行服务
	app.Router.Run(config.CfgServer.HttpPort)
}
