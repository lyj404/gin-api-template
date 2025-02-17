package main

import (
	"gin-api-template/api/route"
	"gin-api-template/bootstrap"
	"gin-api-template/config"
	"gin-api-template/pkg/lib/logger"
	"time"
)

func main() {
	app := bootstrap.App()

	// 初始化自定义日志
	logger := logger.InitLogger()

	// 程序结束时关闭MySQL和Redis连接
	defer app.CloseConnection()

	// 设置超时时间
	timeout := time.Duration(config.CfgTimeout.ContextTimeout) * time.Second

	// 设置swagger文档
	config.SetUpSwag()

	// 设置路由
	router := route.SetUp(timeout, app, logger)
	// 运行web服务
	router.Run(config.CfgServer.HttpPort)
}
