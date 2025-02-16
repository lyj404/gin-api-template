package config

import "gin-api-template/docs"

func SetUpSwag() {
	docs.SwaggerInfo.Title = "gin-api-template"
	docs.SwaggerInfo.Description = "gin-api-template 相关API文档"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "127.0.0.1" + CfgServer.HttpPort
	docs.SwaggerInfo.BasePath = "/"
}
