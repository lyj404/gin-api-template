package main

import (
	"gin-api-template/api/route"
	"gin-api-template/bootstrap"
	"gin-api-template/config"
	"gin-api-template/pkg/lib/logger"
	"time"
)

func main() {
	// 初始化数据库和缓存
	bootstrap.Boot()

	// 初始化自定义日志
	logger := logger.InitLogger()

	// 程序结束时关闭MySQL和Redis连接
	defer bootstrap.CloseConnection()

	// 设置超时时间
	timeout := time.Duration(config.CfgTimeout.ContextTimeout) * time.Second

	// 设置swagger文档
	config.SetUpSwag()

	// 设置路由并运行服务
	route.SetUp(timeout, logger).Run(config.CfgServer.HttpPort)
}
