package main

import (
	"gin-api-template/api/middleware"
	"gin-api-template/api/route"
	"gin-api-template/bootstrap"
	"gin-api-template/config"
	"gin-api-template/domain"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	app := bootstrap.App()

	// 初始化自定义日志
	logger := domain.InitLogger()

	// 程序结束时关闭MySQL和Redis连接
	defer app.CloseConnection()

	// 设置超时时间
	timeout := time.Duration(config.CfgTimeout.ContextTimeout) * time.Second

	gin := gin.New()
	// 使用自定义日志
	gin.Use(middleware.LoggerMiddleware(logger))

	route.SetUp(timeout, app, gin)

	gin.Run(config.CfgServer.HttpPort)
}
