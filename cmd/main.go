package main

import (
	"time"

	"github.com/lyj404/gin-api-template/api/route"
	"github.com/lyj404/gin-api-template/bootstrap"
	"github.com/lyj404/gin-api-template/config"
	"github.com/lyj404/gin-api-template/pkg/lib/logger"
)

func main() {
	// 初始化配置
	config.InitConfig()

	// 初始化数据库和缓存
	bootstrap.Boot()

	// 初始化zap日志
	logger := logger.InitZapLogger()
	defer logger.Sync()

	// 程序结束时关闭MySQL和Redis连接
	defer bootstrap.CloseConnection()

	// 设置超时时间
	timeout := time.Duration(config.CfgTimeout.ContextTimeout) * time.Second

	// 设置swagger文档
	config.SetUpSwag()

	// 设置路由并运行服务
	route.SetUp(timeout, logger).Run(config.CfgServer.HttpPort)
}
