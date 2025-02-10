package main

import (
	"gin-api-template/api/route"
	"gin-api-template/bootstrap"
	"gin-api-template/config"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	app := bootstrap.App()

	// 程序结束时关闭MySQL和Redis连接
	defer app.CloseConnection()

	// 设置超时时间
	timeout := time.Duration(config.ContextTimeOut) * time.Second

	gin := gin.Default()

	route.SetUp(timeout, app, gin)

	gin.Run(config.HttpPort)
}
